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
	"tbank/bot/internal/closer"
	grpcserver "tbank/bot/internal/grpcserver"
	kafkaconsumer "tbank/bot/internal/kafka_consumer"
	telegramproducer "tbank/bot/internal/telegram_producer"
	"tbank/bot/internal/usecase"
	"time"
	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"gopkg.in/telebot.v3"
)


type App struct {
	grpcServer 		*grpc.Server
	config 			*config.Config
	bot				*telebot.Bot
	tgProducer 		*telegramproducer.TelegramProducer
	kafkaConsumer 	*kafkaconsumer.KafkaConsumer
	closer			*closer.Closer
}


func NewApp() (*App, error) {

	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	logger := slog.Default()

	closer := closer.NewCloser(logger)

	grpcServer := grpc.NewServer()

	closer.Add(func() error {
		grpcServer.Stop()
		return nil
	})

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

	closer.Add(func() error {
		bot.Stop()
		return nil
	})


	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(config.Kafka.Addresses, saramaConfig)
	if err != nil {
		return nil, err
	}

	commitCh := make(chan sarama.ConsumerMessage)

	telegramProducer := telegramproducer.NewTelegramProducer(bot, logger, commitCh)

	closer.Add(func() error {
		telegramProducer.Stop()
		return nil
	})

	kafkaConsumer := kafkaconsumer.NewKafkaConsumer(consumer, config.Kafka.Topic, commitCh)

	closer.Add(func() error {
		kafkaConsumer.Stop()
		return nil
	})

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
		tgProducer: telegramProducer,
		kafkaConsumer: kafkaConsumer,
		closer: closer,
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


	a.kafkaConsumer.Run()

	a.tgProducer.Run()

	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <- errorCh:
		if errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("grpc server stopped", slog.String("err", err.Error()))
			return a.closer.Close()
		}
		slog.Error("failed to start the grpc-server")
		return a.closer.Close()
	case <- quit:
		slog.Info("Shutting down")
		return a.closer.Close()
	}
}