package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tbank/scrapper/config"
	grpcserver "tbank/scrapper/internal/grpc-server"
	"tbank/scrapper/api/proto/gen"
	"google.golang.org/grpc"
)

type App struct {
	grpcServer *grpc.Server
	config     *config.Config
	logger     *slog.Logger
}

func NewApp(config *config.Config) *App {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	grpcServer := grpc.NewServer()
	scrapperGRPCServer := grpcserver.NewScrapperServer()
	gen.RegisterScrapperServer(grpcServer, scrapperGRPCServer)

	return &App{
		grpcServer: grpcServer,
		config:     config,
		logger:     logger,
	}
}

func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%s", a.config.ScrapperServer.Host, a.config.ScrapperServer.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	a.logger.Info("Запуск gRPC сервера", "addr", addr)

	errorCh := make(chan error, 1)

	go func() {
		if err := a.grpcServer.Serve(listener); err != nil {
			errorCh <- fmt.Errorf("gRPC server error: %w", err)
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
		a.logger.Error("Ошибка сервера", "error", err)
		return err
	}
}
