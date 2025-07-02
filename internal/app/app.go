package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"order-service/internal/config"
	"order-service/internal/handlers"
	kafkaconsumer "order-service/internal/infrastructure/consumer/kafka"
	kafkaproducer "order-service/internal/infrastructure/producer/kafka"
	projectionworker "order-service/internal/infrastructure/projection-worker"
	"order-service/internal/infrastructure/worker"
	"order-service/internal/repository/cache/redis"
	"order-service/internal/repository/storage/postgres"
	"order-service/internal/usecase"
	"order-service/pkg/logger"
	"order-service/pkg/txmanager"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)




type App struct {
	server 			 *http.Server
	logger 			 logger.Logger
	worker 			 *worker.Worker
	projectionWorker *projectionworker.ProjectionWorker
}


func NewApp() (*App, error) {

	cfg, err := config.LoadConfig(".env", "./configs")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	router := gin.Default()

	gin.SetMode(gin.ReleaseMode)

	logger, err := SetupLogger(config.LoggingConfig{
		Level:             "info",
		Mode:              "dev",
		Encoding:          "console",
		DisableCaller:     true,
		DisableStacktrace: true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		TimestampKey: 		"timestamp",
		CapitalizeLevel:    true,
		InitialFields: map[string]interface{}{
			"service": "order-service",
			"env":     "productio",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	mainPool, err := cfg.Storage.Main.Connect(context.Background())
	if err != nil {
		logger.Errorw("Failed to connect to main database", "error", err)
		return nil, err
	}

	if err := cfg.Storage.Main.ApplyMigrations(context.Background(), "./migrations"); err != nil {
		logger.Errorw("Failed to apply migrations to main database", "error", err)
		return nil, err
	}

	logger.Infow("Connected to databases", "mainDB", cfg.Storage.Main.Host, "sideDB", cfg.Storage.Side.Host)


	main_txmanager 	 	:= txmanager.NewPgxTxManager(mainPool)

	mainOrderStorage 	:= postgres.NewPostgresOrderStorage(main_txmanager)
	outboxStorage 	 	:= postgres.NewOutboxStorage(main_txmanager)

	redisCache 	     	:= redis.NewRedisCache(cfg.Redis.Address)

	createOrderUseCase 	:= usecase.NewCreateOrderUseCase(mainOrderStorage, outboxStorage, main_txmanager, logger)
	getOrdersUseCase 	:= usecase.NewGetOrdersUseCase(redisCache)
	deleteOrderUseCase 	:= usecase.NewDeleteOrderUseCase(mainOrderStorage, outboxStorage, main_txmanager, logger)

	handler 			:= handlers.NewHandler(createOrderUseCase, deleteOrderUseCase, getOrdersUseCase, logger)

	consumer, err 		:= kafkaconsumer.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	producer 			:= kafkaproducer.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic) 

	worker 				:= worker.NewWorker(producer, consumer, outboxStorage, main_txmanager, logger)

	projectionWorker 	:= projectionworker.NewProjectionWorker(consumer, redisCache, logger)

	SetRoutes(router, handler)
	
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
	}

	return &App{
		server: 			 server,
		logger: 			 logger,
		worker: 			 worker,
		projectionWorker: 	 projectionWorker,
	}, nil
}

func (a *App) Start() error {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Infow("Starting server", "address", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
			return
		}
	}()

	go func() {
		a.logger.Infow("Starting worker")
		if err := a.worker.Run(context.Background()); err != nil && errors.Is(err, context.Canceled) {
			a.logger.Warnw("Worker stopped", "error", err)
			return
		}
	}()

	go func() {
		a.logger.Infow("Starting projection worker")
		if err := a.projectionWorker.Run(); err != nil  && errors.Is(err, context.Canceled) {
			a.logger.Warnw("Projection worker stopped", "error", err)
			return
		}
	}()

	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.logger.Infow("Shutting down server gracefully")
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	a.logger.Infow("Worker stopping")
	a.worker.Stop()

	a.logger.Infow("Projection worker stopping")
	a.projectionWorker.Stop()

	time.Sleep(1 * time.Second)

	return nil
}



