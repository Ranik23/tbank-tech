package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tbank/scrapper/config"
	"tbank/scrapper/internal/router"
	"time"
)

type App struct {
	server *http.Server
	config *config.Config
	logger *slog.Logger
}

func NewApp(config *config.Config) *App {

	addr := fmt.Sprintf("%s:%s", config.ScrapperServer.Host, config.ScrapperServer.Port)

	router := router.NewRouter()

	// usecase

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &App{
		server: srv,
		config: config,
	}
}

func (a *App) Run() error {

	errorCh := make(chan error, 1)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			errorCh <- err
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		a.logger.Info("Получен сигнал завершения, выключаем сервер...")
		if err := a.server.Shutdown(ctx); err != nil {
			return err
		}

		a.logger.Error("Сервер выключен корректно")
		return nil
	case err := <-errorCh:
		a.logger.Error("Ошибка сервера")
		return err
	}
}
