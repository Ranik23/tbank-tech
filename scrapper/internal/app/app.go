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
	"tbank/scrapper/internal/gateway"
	grpcserver "tbank/scrapper/internal/controllers/grpc"
	"tbank/scrapper/internal/hub"
	git "tbank/scrapper/pkg/github"

	"tbank/scrapper/internal/usecase"

	// "github.com/IBM/sarama"
	"github.com/google/go-github/v69/github"
	"google.golang.org/grpc"
)

type App struct {
	grpcServer *grpc.Server
	config     *config.Config
	logger     *slog.Logger
}

func NewApp() (*App, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	grpcServer := grpc.NewServer()

	// storage, err := storage.NewStorageImpl(cfg)
	// if err != nil {
	// 	return nil, err
	// }

	// producer, err := sarama.NewAsyncProducer(cfg.Kafka.Addresses, nil)
	// if err != nil {
	// 	return nil, err
	// }

	gitHubClient := git.NewRealGitHubClient()

	commitCh := make(chan *github.RepositoryCommit)

	hub := hub.NewHub(gitHubClient, commitCh, slog.Default())

	usecase , err := usecase.NewUseCaseImpl(cfg, nil, hub, logger)
	if err != nil {
		return nil, err
	}

	grpcScrapperServer := grpcserver.NewScrapperServer(usecase)
	

	gen.RegisterScrapperServer(grpcServer, grpcScrapperServer)

	return &App{
		grpcServer: grpcServer,
		config:     cfg,
		logger:     logger,
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
		a.grpcServer.GracefulStop()
		a.logger.Info("gRPC сервер корректно завершил работу")
		return nil
	case err := <-errorCh:
		a.logger.Error("Ошибка сервера", slog.String("error", err.Error()))
		return err

	case err := <- errorChProxy:
		a.logger.Error("Ошибка прокси-сервера", slog.String("error", err.Error()))
		return err
	}
}
