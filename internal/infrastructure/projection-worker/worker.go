package projectionworker

import (
	"context"
	"errors"
	"order-service/internal/config"
	"order-service/internal/domain"
	"order-service/internal/infrastructure/consumer"
	"order-service/internal/infrastructure/consumer/kafka"
	"order-service/internal/repository/cache"
	"order-service/pkg/logger"
	"time"

	"golang.org/x/sync/errgroup"
)

type ProjectionWorker struct {
	cacheClient 		cache.Cache
	logger 				logger.Logger
	cfg                 *config.Config


	ctx 			  context.Context
	cancel            context.CancelFunc
}

func NewProjectionWorker(cache cache.Cache, logger logger.Logger, config *config.Config) *ProjectionWorker {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &ProjectionWorker{
		cacheClient:    cache,
		logger:   logger,
		cfg:      config,
		ctx:     ctx,
		cancel:  cancelFunc,
	}
}

func (w *ProjectionWorker) Run() error {
	w.logger.Infow("Projection worker started")
	defer w.logger.Infow("Projection worker stopped")

	g, errctx := errgroup.WithContext(context.Background())

	for workerID := 0; workerID < w.cfg.Kafka.NumWorkers; workerID++ {
		g.Go(func() error {
			consumer, err := kafka.NewKafkaConsumer(w.cfg.Kafka.Brokers, w.cfg.Kafka.Topic, w.cfg.Kafka.GroupID)
			if err != nil {
				w.logger.Errorw("Failed to create Kafka consumer", "error", err)
				return err
			}
			w.logger.Infow("Starting projection worker", "workerID", workerID)
			if err := w.dispatchWorker(consumer, w.ctx, errctx, workerID); err != nil && !errors.Is(err, context.Canceled) {
				w.logger.Errorw("Projection worker stopped with error", "error", err)
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		w.logger.Errorw("Projection worker encountered an error", "error", err)
		return err
	}
	
	return nil
}

func (w *ProjectionWorker) Stop() {
	w.logger.Info("Projection worker stopping")
	w.cancel()
}


func (w *ProjectionWorker) dispatchWorker(consumer consumer.Consumer, ctx context.Context, errctx context.Context, workerID int) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	w.logger.Infow("Projection Worker started", "workerID", workerID)

	for {
		select {
		case <-ctx.Done():
			w.logger.Infow("Projection Worker stopping due to ctx done", "workerID", workerID)
			return context.Canceled
		case <-errctx.Done():
			w.logger.Infow("Projection Worker stopping due to errgroup ctx done", "workerID", workerID)
			return errctx.Err()
		case <-ticker.C:
			if err := w.dispatchEvent(consumer, ctx); err != nil {
				w.logger.Errorw("Failed to dispatch event", "error", err, "workerID", workerID)
				continue
			}
		}
	}
}


func (w *ProjectionWorker) dispatchEvent(consumer consumer.Consumer, ctx context.Context) error {
	value, key, err := consumer.Consume(ctx)
	if err != nil {
		w.logger.Errorw("Failed to consume message", "error", err)
		return err
	}

	switch string(key) {
	case domain.OrderCreatedKey:
		if err := w.handleCreateOrder(ctx, value); err != nil {
			w.logger.Errorw("Failed to handle create order", "error", err)
			return err
		}
		w.logger.Infow("Order sent to Redis successfully", "order", string(value))
	case domain.OrderDeletedKey:
		if err := w.handleDeleteOrder(ctx, value); err != nil {
			w.logger.Errorw("Failed to handle delete order", "error", err)
			return err
		}
		w.logger.Infow("Order deleted from Redis successfully", "order", string(value))
	default:
		w.logger.Errorw("Unknown message key", "key", string(key))
		return errors.New("unknown message key: " + string(key))
	}

	return nil
}
