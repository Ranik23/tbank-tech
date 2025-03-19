package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ranik23/tbank-tech/scrapper/api/proto/gen"
	"github.com/Ranik23/tbank-tech/scrapper/config"
	grpcserver "github.com/Ranik23/tbank-tech/scrapper/internal/controllers/grpc"
	"github.com/Ranik23/tbank-tech/scrapper/internal/gateway"
	"github.com/Ranik23/tbank-tech/scrapper/internal/hub"
	kafkaproducer "github.com/Ranik23/tbank-tech/scrapper/internal/kafka_producer"
	"github.com/Ranik23/tbank-tech/scrapper/internal/repository/postgres"
	"github.com/Ranik23/tbank-tech/scrapper/internal/service"
	git "github.com/Ranik23/tbank-tech/scrapper/pkg/github_client"

	"github.com/IBM/sarama"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
)

type App struct {
	grpcServer    *grpc.Server
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
		logger.Error("Failed to laod the config", slog.String("error", err.Error()))
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

	grpcServer := grpc.NewServer()

	closer.Add(func(ctx context.Context) error {
		logger.Info("Stopping gRPC serve...")
		grpcServer.GracefulStop()
		return nil
	})

	gitHubClient := git.NewRealGitHubClient(logger)

	commitCh := make(chan hub.CustomCommit)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(cfg.Kafka.Addresses, saramaConfig)
	if err != nil {
		logger.Error("Failed to create a new async Kafka producer", slog.String("error", err.Error()))
		return nil, err
	}

	kafkaProducer, err := kafkaproducer.NewKafkaProducer(producer, logger, commitCh, cfg.Kafka.Topic)
	if err != nil {
		logger.Error("Failed to create a new Kafka producer", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Successfully created a Kafka producer")

	closer.Add(func(ctx context.Context) error {
		logger.Info("Stopping Kafka Producer...")
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
	}, nil
}

func (a *App) Run() error {

	grpcAddr := fmt.Sprintf("%s:%s", a.config.ScrapperServer.Host, a.config.ScrapperServer.Port)
	httpAddr := fmt.Sprintf("%s:%s", a.config.ScrapperServerHTTP.Host, a.config.ScrapperServerHTTP.Port)

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	a.logger.Info("Запуск gRPC сервера", "grpcAddr", grpcAddr)

	errorCh := make(chan error, 2)

	a.hub.Run()

	a.kafkaProducer.Run()

	go func() {
		if err := a.grpcServer.Serve(listener); err != nil {
			errorCh <- fmt.Errorf("gRPC server error: %v", err)
		}
	}()

	go func() {
		if err := gateway.RunGateway(context.Background(), grpcAddr, httpAddr, a.logger); err != nil {
			errorCh <- fmt.Errorf("http proxy server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		a.logger.Info("Получен сигнал завершения, выключаем gRPC сервер...")
		if err := a.closer.Close(context.Background()); err != nil {
			a.logger.Error("Ошибка при закрытии ресурсов", slog.String("error", err.Error()))
		} else {
			a.logger.Info("Все ресурсы корректно закрыты")
		}
		a.logger.Info("gRPC сервер корректно завершил работу")
		return nil
	case err := <-errorCh:
		a.logger.Error("Ошибка сервера", slog.String("error", err.Error()))
		if err := a.closer.Close(context.Background()); err != nil {
			a.logger.Error("Ошибка при закрытии ресурсов", slog.String("error", err.Error()))
		}
		return err
	}
}
