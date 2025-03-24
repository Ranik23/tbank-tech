package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ranik23/tbank-tech/scrapper/api/proto/gen"
	"github.com/Ranik23/tbank-tech/scrapper/config"
	grpcserver "github.com/Ranik23/tbank-tech/scrapper/internal/controllers/grpc"
	"github.com/Ranik23/tbank-tech/scrapper/internal/gateway"
	"github.com/Ranik23/tbank-tech/scrapper/internal/hub"
	kafkaproducer "github.com/Ranik23/tbank-tech/scrapper/internal/kafka"
	"github.com/Ranik23/tbank-tech/scrapper/internal/metrics"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
	"github.com/Ranik23/tbank-tech/scrapper/internal/service"
	git "github.com/Ranik23/tbank-tech/scrapper/pkg/github_client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/IBM/sarama"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
)

func init() {
	prometheus.MustRegister(metrics.TotalRequests, metrics.ErrorRequests, metrics.RequestDuration)
}

type App struct {
	grpcServer    *grpc.Server
	metricsServer *http.Server
	config        *config.Config
	logger        *slog.Logger
	kafkaProducer *kafkaproducer.KafkaProducer
	hub           hub.Hub
	closer        *Closer
}

func NewApp() (*App, error) {

	logger := slog.New(tint.NewHandler(os.Stdout, nil)).With(slog.String("SERVICE", "SCRAPPER"))

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Error("Failed to load the config", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Successfully loaded the config")

	pool, err := cfg.DataBase.ConnectToPostgres(context.Background())
	if err != nil {
		logger.Error("Failed to connect to DB", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Successfully connected to PostgreSQL")

	closer := NewCloser()

	closer.Add(func(ctx context.Context) error {
		logger.Info("Closing pgx.Pool...")
		pool.Close()
		return nil
	})

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpcserver.ErrorHandlingInterceptor))

	closer.Add(func(ctx context.Context) error {
		logger.Info("Stopping gRPC serve...")
		grpcServer.GracefulStop()
		return nil
	})

	metricAddr := fmt.Sprintf("%s:%s", cfg.MetricServer.Host, cfg.MetricServer.Port)

	metricsServer := &http.Server{
		Addr: metricAddr,
	}

	closer.Add(func(ctx context.Context) error {
		logger.Info("Shutting down the METRICS Server!")
		return metricsServer.Shutdown(ctx)
	})

	gitHubClient := git.NewRealGitHubClient(logger)

	commitCh := make(chan hub.CustomCommit)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	kafkaProducer, err := kafkaproducer.NewKafkaProducer(cfg.Kafka.Addresses, logger, commitCh, cfg.Kafka.Topic, saramaConfig)
	if err != nil {
		logger.Error("Failed to create a new Kafka producer", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Successfully created a Kafka producer")

	closer.Add(func(ctx context.Context) error {
		logger.Info("Stopping Kafka Producer...!")
		kafkaProducer.Stop()
		return nil
	})

	hub := hub.NewHub(gitHubClient, commitCh, logger)

	closer.Add(func(ctx context.Context) error {
		logger.Info("Stopping Hub...")
		hub.Stop()
		return nil
	})

	txManager := postgres.NewTxManager(pool, logger)

	postgresRepo := postgres.NewPostgresRepository(txManager, logger)

	usecase, err := service.NewService(postgresRepo, txManager, hub, logger)
	if err != nil {
		logger.Error("Failed to create a new service", slog.String("error", err.Error()))
		return nil, err
	}

	grpcScrapperServer := grpcserver.NewScrapperServer(usecase)

	gen.RegisterScrapperServer(grpcServer, grpcScrapperServer)

	return &App{
		grpcServer:    grpcServer,
		config:        cfg,
		logger:        logger,
		kafkaProducer: kafkaProducer,
		hub:           hub,
		closer:        closer,
		metricsServer: metricsServer,
	}, nil
}

func (a *App) Run() error {

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		if err := a.closer.Close(ctx); err != nil {
			a.logger.Error("Failed to close resources", slog.String("error", err.Error()))
		}
		a.logger.Info("Successfully closed all resources")
	}()

	grpcAddr := fmt.Sprintf("%s:%s", a.config.ScrapperServer.Host, a.config.ScrapperServer.Port)
	httpAddr := fmt.Sprintf("%s:%s", a.config.ScrapperServerHTTP.Host, a.config.ScrapperServerHTTP.Port)
	metricsAddr := fmt.Sprintf("%s:%s", a.config.MetricServer.Host, a.config.MetricServer.Port)

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	errorCh := make(chan error, 2)

	a.hub.Run()

	a.kafkaProducer.Run()

	go func() {
		a.logger.Info("Запуск gRPC сервера", slog.String("grpcAddr", grpcAddr))
		if err := a.grpcServer.Serve(listener); err != nil {
			errorCh <- fmt.Errorf("gRPC server error: %v", err)
		}
	}()

	go func() {
		a.logger.Info("Запукс прокси-сервера", slog.String("httpAddr", httpAddr))
		if err := gateway.RunGateway(context.Background(), grpcAddr, httpAddr, a.logger); err != nil {
			errorCh <- fmt.Errorf("http proxy server error: %v", err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		a.logger.Info("Запуск Prometheus-метрик", slog.String("addr", metricsAddr))
		if err := a.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorCh <- fmt.Errorf("metrics server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		a.logger.Info("Получен сигнал завершения, выключаем gRPC сервер...")
		return nil
	case err := <-errorCh:
		a.logger.Error("Ошибка сервера", slog.String("error", err.Error()))
		return err
	}
}
