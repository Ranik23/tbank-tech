package app

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tbank/bot/config"
	handlers "tbank/bot/internal/bot-handlers"
	"tbank/bot/internal/bot-usecase"
	grpcserver "tbank/bot/internal/grpc-server"
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

	botUseCase := botusecase.NewUseCaseImpl(config, nil, nil)

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

	var users sync.Map

	bot.Handle("/start", handlers.StartHandler(botUseCase, &users))
	bot.Handle("/help", handlers.HelpHandler(botUseCase, &users))
	bot.Handle("/track", handlers.TrackHandler(botUseCase, &users))
	bot.Handle("/untrack", handlers.UnTrackHandler(botUseCase, &users))
	bot.Handle("/list", handlers.ListHandler(botUseCase, &users))
	bot.Handle(telebot.OnText, handlers.MessageHandler(botUseCase, &users))


	grpcBotServer := grpcserver.NewBotServer(botUseCase, bot)

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
			slog.Error("grpc server stopped: %v", err.Error())
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