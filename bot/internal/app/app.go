package app

import (
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tbank/bot/api/proto/gen"
	"tbank/bot/config"
	bothandlers "tbank/bot/internal/bot_handlers"
	botusecase "tbank/bot/internal/bot_usecase"
	grpcserver "tbank/bot/internal/grpcserver"
//	kafkacosumer "tbank/bot/internal/kafka-cosumer"
//	telegramproducer "tbank/bot/internal/telegram-producer"
	"tbank/bot/internal/usecase"
	"time"

	//"github.com/confluentinc/confluent-kafka-go/kafka"
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

	logger := slog.Default()

	grpcServer := grpc.NewServer()

	botUseCase, err := botusecase.NewUseCaseImpl(config, nil, logger)
	if err != nil {
		return nil, err
	}

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

	//messagesCh := make(chan kafka.Message)

	//telegramProducer := telegramproducer.NewTelegramProducer(bot, messagesCh)

	useCase := usecase.NewUseCaseImp(bot)

	var users sync.Map

	bot.Handle("/start", bothandlers.StartHandler(botUseCase, &users))
	bot.Handle("/help", bothandlers.HelpHandler(botUseCase, &users))
	bot.Handle("/track", bothandlers.TrackHandler(botUseCase, &users))
	bot.Handle("/untrack", bothandlers.UnTrackHandler(botUseCase, &users))
	bot.Handle("/list", bothandlers.ListHandler(botUseCase, &users))
	bot.Handle(telebot.OnText, bothandlers.MessageHandler(botUseCase, &users))

	grpcBotServer := grpcserver.NewBotServer(useCase, bot)

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
			slog.Error("grpc server stopped", slog.String("err", err.Error()))
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