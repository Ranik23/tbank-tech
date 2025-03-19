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
	"time"

	"github.com/Ranik23/tbank-tech/bot/api/proto/gen"
	"github.com/Ranik23/tbank-tech/bot/config"
	grpcserver "github.com/Ranik23/tbank-tech/bot/internal/controllers/grpc"
	telegramhandlers "github.com/Ranik23/tbank-tech/bot/internal/controllers/telegram"
	kafkaconsumer "github.com/Ranik23/tbank-tech/bot/internal/kafka_consumer"
	"github.com/Ranik23/tbank-tech/bot/internal/service"
	telegramproducer "github.com/Ranik23/tbank-tech/bot/internal/telegram_producer"

	"github.com/IBM/sarama"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	logger := slog.New(tint.NewHandler(os.Stdout, nil))

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load the config", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Successfully loaded the config")

	closer := NewCloser()

	grpcServer := grpc.NewServer()

	closer.Add(func(ctx context.Context) error {
		grpcServer.GracefulStop()
		return nil
	})

	connectionStr := fmt.Sprintf("%s:%s", config.ScrapperService.Host, config.ScrapperService.Port)

	connScrapper, err  := grpc.NewClient(connectionStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to establish connection with Scrapper Service", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Info("Successfully connected to gRPC Scrapper Server")

	closer.Add(func(ctx context.Context) error {
		if err := connScrapper.Close(); err != nil {
			return err
		}
		return nil
	})

	botService := service.NewService(connScrapper, config, logger)


	bot, err := telebot.NewBot(telebot.Settings{
		Token: config.Telegram.Token,
		Poller: &telebot.LongPoller {
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
	logger.Info("Successfully created a Kafka consumer")



	commitCh := make(chan sarama.ConsumerMessage)

	telegramProducer := telegramproducer.NewTelegramProducer(bot, logger, commitCh)

	closer.Add(func(ctx context.Context) error {
		telegramProducer.Stop()
		return nil
	})

	kafkaConsumer := kafkaconsumer.NewKafkaConsumer(consumer, config.Kafka.Topic, commitCh, logger)

	closer.Add(func(ctx context.Context) error {
		kafkaConsumer.Stop()
		return nil
	})

	var users sync.Map

	botHandler := telegramhandlers.NewBotHandler(botService, &users)

	bot.Handle("/start", botHandler.StartHandler())
	bot.Handle("/help", botHandler.HelpHandler())
	bot.Handle("/track", botHandler.TrackHandler())
	bot.Handle("/untrack", botHandler.UnTrackHandler())
	bot.Handle("/list", botHandler.ListLinksHandler())
	bot.Handle(telebot.OnText, botHandler.MessageHandler())

	grpcBotServer := grpcserver.NewBotGRPCServer(bot, botService)

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

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
		if err := a.closer.Close(ctx); err != nil {
			a.logger.Error("Failed to close resources", slog.String("error", err.Error()))
		}
	}()

	a.logger.Info("Starting the bot...")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", a.config.TelegramBotServer.Host, a.config.TelegramBotServer.Port))
	if err != nil {
		return err
	}

	errorCh := make(chan error, 1)

	go func() {
		a.logger.Info("Starting gRPC server", slog.String("port", a.config.TelegramBotServer.Port))
		if err := a.grpcServer.Serve(listener); err != nil {
			errorCh <- err
		}
	}()

	go func() {
		a.logger.Info("Starting Telegram bot...")
		a.bot.Start()
	}()

	
	go func() {
		a.logger.Info("Starting Kafka Consumer...")
		if err := a.kafkaConsumer.Run(); err != nil {
			errorCh <- err
		}
	}()

	go func() {
		a.tgProducer.Run()
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <- errorCh:
		if errors.Is(err, grpc.ErrServerStopped) {
			a.logger.Error("gRPC server stopped", slog.String("err", err.Error()))
			return nil
		}
		a.logger.Error("Service Failed", slog.String("error", err.Error()))
		return err
	case <- quit:
		a.logger.Info("Shutting down")
		return nil
	}
}