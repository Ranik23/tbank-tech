package app

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tbank/bot/config"
	"tbank/bot/internal/server"
	"tbank/bot/internal/server/handlers"
	"tbank/bot/internal/server/router"
	"tbank/bot/internal/usecase"
	"tbank/bot/proto/gen"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/telebot.v3"
)


type App struct {
	grpcServer *grpc.Server
	config 		*config.Config
	bot			*telebot.Bot
}


func NewApp() (*App, error) {

	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	usecase := usecase.NewUseCaseImpl(config, nil, nil)

	bot, err := telebot.NewBot(telebot.Settings{
		Token: config.Telegram.Token,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
			AllowedUpdates: []string{
				"message",
				"edited_message",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	botRouter := router.NewRouter(bot, nil)

	botRouter.AddHandler("/start", handlers.NewStartHandler(usecase))
	botRouter.AddHandler("/help", handlers.NewHelpHandler(usecase))
	botRouter.AddHandler("/track", handlers.NewTrackHandler(usecase))
	botRouter.AddHandler("/untrack", handlers.NewUntrackHandler(usecase))
	botRouter.AddHandler("/list", handlers.NewListHandler(usecase))

	botRouter.RegisterHandlers()

	grpcBotServer := server.NewBotServer(usecase, bot)

	gen.RegisterBotServer(grpcServer, grpcBotServer)

	return &App{
		grpcServer: grpcServer,
		config: config,
		bot: bot,
	}, nil
}


func (a *App) Run() error {

	listener, err := net.Listen("tcp", ":" + a.config.TelegramBotServer.Port)
	if err != nil {
		return err
	}
	
	errorCh := make(chan error, 1)

	go func() {
		slog.Info("Starting gRPC server on port " + a.config.TelegramBotServer.Port)
		if err := a.grpcServer.Serve(listener); err != nil {
			errorCh <- err
		}
	}()

	go func() {
		slog.Info("Starting Telegram bot...")
		a.bot.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	

	select {
	case err := <- errorCh:
		if errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("grpc server stopped: %v", err)
			return nil
		}
		slog.Error("failed to start the grpc-server")
		return err
	case <- quit:
		slog.Info("Shutting down gRPC server...")
		a.grpcServer.GracefulStop()
		return nil
	}
}