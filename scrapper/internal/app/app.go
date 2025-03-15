package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tbank/scrapper/api/proto/gen"
	"tbank/scrapper/config"
	"tbank/scrapper/internal/closer"
	grpcserver "tbank/scrapper/internal/controllers/grpc"
	"tbank/scrapper/internal/gateway"
	"tbank/scrapper/internal/hub"
	kafkaproducer "tbank/scrapper/internal/kafka_producer"
	"tbank/scrapper/internal/service"
	git "tbank/scrapper/pkg/github"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
)


type App struct {
	grpcServer 		*grpc.Server
	config     		*config.Config
	logger     		*slog.Logger
	kafkaProducer 	*kafkaproducer.KafkaProducer
	hub				hub.Hub
	closer			*closer.Closer
}

func NewApp() (*App, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	closer := closer.NewCloser(logger)

	grpcServer := grpc.NewServer()

	closer.Add(func() error {
		grpcServer.GracefulStop()
		return nil
	})

	gitHubClient := git.NewRealGitHubClient()

	commitCh := make(chan hub.CustomCommit)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer([]string{"localhost:9093"}, saramaConfig) 
	if err != nil {
		return nil, err
	}

	kafkaTopic := cfg.Kafka.Topic
	
	kafkaProducer, err := kafkaproducer.NewKafkaProducer(producer, logger, commitCh, kafkaTopic)
	if err != nil {
		return nil, err
	}

	closer.Add(func() error {
		kafkaProducer.Stop()
		return nil
	})

	hub := hub.NewHub(gitHubClient, commitCh, slog.Default())

	closer.Add(func() error {
		hub.Stop()
		return nil
	})

	usecase , err := service.NewService(nil, hub, logger)
	if err != nil {
		return nil, err
	}

	grpcScrapperServer := grpcserver.NewScrapperServer(usecase)
	
	gen.RegisterScrapperServer(grpcServer, grpcScrapperServer)

	return &App{
		grpcServer: 	grpcServer,
		config:     	cfg,
		logger:     	logger,
		kafkaProducer:	kafkaProducer,
		hub: 			hub,
		closer: 		closer,
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

	errorCh := make(chan error, 1)

	errorChProxy := make(chan error, 1)

	a.hub.Run()

	a.kafkaProducer.Run()

	go func() {
		if err := a.grpcServer.Serve(listener); err != nil {
			errorCh <- fmt.Errorf("gRPC server error: %v", err)
		}
	}()

	go func() {
		if err := gateway.RunGateway(context.Background(), grpcAddr, httpAddr); err != nil {
			errorChProxy <- fmt.Errorf("http proxy server error: %v", err)
		}
	}()


	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		a.logger.Info("Получен сигнал завершения, выключаем gRPC сервер...")
		if err := a.closer.Close(); err != nil {
			a.logger.Error("Ошибка при закрытии ресурсов", slog.String("error", err.Error()))
		} else {
			a.logger.Info("Все ресурсы корректно закрыты")
		}
		a.logger.Info("gRPC сервер корректно завершил работу")
		return nil
	case err := <-errorCh:
		a.logger.Error("Ошибка сервера", slog.String("error", err.Error()))
		if err := a.closer.Close(); err != nil {
			a.logger.Error("Ошибка при закрытии ресурсов", slog.String("error", err.Error()))
		}
		return err
	case err := <-errorChProxy:
		a.logger.Error("Ошибка прокси-сервера", slog.String("error", err.Error()))
		if err := a.closer.Close(); err != nil {
			a.logger.Error("Ошибка при закрытии ресурсов", slog.String("error", err.Error()))
		}
		return err
	}
}
