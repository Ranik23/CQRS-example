package projectionworker

import (
	"context"
	"errors"
	"order-service/internal/domain"
	"order-service/internal/infrastructure/consumer"
	"order-service/internal/repository/cache"
	"order-service/pkg/logger"
	"time"
)

type ProjectionWorker struct {
	consumer    		consumer.Consumer
	cacheClient 		cache.Cache
	logger 				logger.Logger


	ctx 			  context.Context
	cancel            context.CancelFunc
}

func NewProjectionWorker(consumer consumer.Consumer, cache cache.Cache, logger logger.Logger) *ProjectionWorker {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &ProjectionWorker{
		consumer: consumer,
		cacheClient:    cache,
		logger:   logger,
		ctx:     ctx,
		cancel:  cancelFunc,
	}
}

func (w *ProjectionWorker) Run() error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer w.cancel()

	for {
		select {
		case <-w.ctx.Done():
			w.logger.Infow("Projection worker stopped")
			return context.Canceled
		case <-ticker.C:
			err := w.dispatch(w.ctx)
			if err != nil {
				w.logger.Errorw("Failed to dispatch message", "error", err)
				continue
			}
		}
	}
}

func (w *ProjectionWorker) Stop() {
	w.logger.Info("Projection worker stopping")
	w.cancel()
}

func (w *ProjectionWorker) dispatch(ctx context.Context) error {
	value, key, err := w.consumer.Consume(ctx)
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
