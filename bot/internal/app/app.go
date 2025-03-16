package app

import (
	"context"
	"errors"
	"fmt"
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
	grpcserver "tbank/bot/internal/controllers/grpc"
	kafkaconsumer "tbank/bot/internal/kafka_consumer"
	telegramproducer "tbank/bot/internal/telegram_producer"
	"time"

	"github.com/IBM/sarama"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
	"gopkg.in/telebot.v3"
)


type App struct {
	grpcServer 		*grpc.Server
	config 			*config.Config
	bot				*telebot.Bot
	tgProducer 		*telegramproducer.TelegramProducer
	kafkaConsumer 	*kafkaconsumer.KafkaConsumer
	closer			*Closer
	logger			*slog.Logger
}


func NewApp() (*App, error) {

	logger := slog.New(tint.NewHandler(os.Stdout, nil)).With(slog.String("SERVICE", "BOT"))

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load the config", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Successfully loaded the config")

	closer := NewCloser()

	grpcServer := grpc.NewServer()

	closer.Add(func(ctx context.Context) error {
		grpcServer.Stop()
		return nil
	})

	botUseCase, err := botusecase.NewUseCaseImpl(config, logger)
	if err != nil {
		logger.Error("Failed to establish the connection to gRPC Scrapper Server", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Successfully connected to gRPC Scrapper Server")

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
		logger.Error("Failed to initialize the bot", slog.String("error", err.Error()))
		return nil, err	
	}

	closer.Add(func(ctx context.Context) error {
		bot.Stop()
		return nil
	})


	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(config.Kafka.Addresses, saramaConfig)
	if err != nil {
		logger.Error("Failed to create a new Sarama consumer", slog.String("error", err.Error()))
		return nil, err
	}

	commitCh := make(chan sarama.ConsumerMessage)

	telegramProducer := telegramproducer.NewTelegramProducer(bot, logger, commitCh)

	closer.Add(func(ctx context.Context) error {
		telegramProducer.Stop()
		return nil
	})

	kafkaConsumer := kafkaconsumer.NewKafkaConsumer(consumer, config.Kafka.Topic, commitCh, logger)

	logger.Info("Successfully created a Kafka consumer")

	closer.Add(func(ctx context.Context) error {
		kafkaConsumer.Stop()
		return nil
	})

	var users sync.Map

	bot.Handle("/start", bothandlers.StartHandler(botUseCase, &users))
	bot.Handle("/help", bothandlers.HelpHandler(botUseCase, &users))
	bot.Handle("/track", bothandlers.TrackHandler(botUseCase, &users))
	bot.Handle("/untrack", bothandlers.UnTrackHandler(botUseCase, &users))
	bot.Handle("/list", bothandlers.ListHandler(botUseCase, &users))
	bot.Handle(telebot.OnText, bothandlers.MessageHandler(botUseCase, &users))

	grpcBotServer := grpcserver.NewBotServer(bot)

	gen.RegisterBotServer(grpcServer, grpcBotServer)

	return &App{
		grpcServer: grpcServer,
		config: config,
		bot: bot,
		tgProducer: telegramProducer,
		kafkaConsumer: kafkaConsumer,
		closer: closer,
		logger: logger,
	}, nil
}


func (a *App) Run() error {

	a.logger.Info("Starting the bot...")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.config.TelegramBotServer.Host, a.config.TelegramBotServer.Port))
	if err != nil {
		return err
	}
	
	errorCh := make(chan error, 1)

	go func() {
		a.logger.Info("Starting gRPC server on port " + a.config.TelegramBotServer.Port)
		if err := a.grpcServer.Serve(listener); err != nil {
			errorCh <- err
		}
	}()

	go func() {
		a.logger.Info("Starting Telegram bot...")
		a.bot.Start()
	}()


	a.kafkaConsumer.Run()

	a.tgProducer.Run()

	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <- errorCh:
		if errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("grpc server stopped", slog.String("err", err.Error()))
			return a.closer.Close(context.Background())
		}
		slog.Error("failed to start the grpc-server")
		return a.closer.Close(context.Background())
	case <- quit:
		slog.Info("Shutting down")
		return a.closer.Close(context.Background())
	}
}