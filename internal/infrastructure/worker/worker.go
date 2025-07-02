package worker

import (
	"context"
	"errors"
	"order-service/internal/infrastructure/consumer"
	"order-service/internal/infrastructure/producer"
	"order-service/internal/repository/storage"
	"order-service/internal/repository/storage/errs"
	"order-service/pkg/logger"
	"order-service/pkg/txmanager"
	"time"
)

type Worker struct {
	outboxStorage storage.OutboxStorage
	producer      producer.Producer
	consumer      consumer.Consumer
	txmanager     txmanager.TxManager

	logger logger.Logger
	ctx context.Context
	cancel context.CancelFunc
}

func NewWorker(producer producer.Producer, consumer consumer.Consumer, outboxStorage storage.OutboxStorage, txmanager txmanager.TxManager, logger logger.Logger) *Worker {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Worker{
		producer:      producer,
		consumer:      consumer,
		outboxStorage: outboxStorage,
		txmanager:     txmanager,	
		logger:        logger,
		ctx:           ctx,
		cancel:        cancelFunc,
	}
}

func (w *Worker) Run() error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() 
	for {
		select {
		case <-w.ctx.Done():
			return context.Canceled
		case <-ticker.C:
			if err := w.dispatchEvent(w.ctx); err != nil {
				w.logger.Warnw("Failed to produce event", "error", err)
				continue
			}
		}
	}
}

func (w *Worker) Stop() {
	w.cancel()
	w.logger.Infow("Worker stopped")
}

func (w *Worker) dispatchEvent(ctx context.Context) error {
	msg, err := w.outboxStorage.GetOutBoxMessage(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNoFound) {
			w.logger.Debugw("No outbox message found, waiting for next tick")
			return nil
		}
		w.logger.Errorw("Failed to fetch message from outbox", "error", err)
		return err
	}

	if err := w.producer.Produce(ctx, []byte(msg.Key), msg.Message); err != nil {
		w.logger.Errorw("Failed to produce message", "error", err, "messageID", msg.ID)
		return err
	}
	
	if err := w.txmanager.Run(ctx, func(ctx context.Context) error {
		if err := w.outboxStorage.MarkAsSent(ctx, msg.ID); err != nil {
			w.logger.Errorw("Failed to mark message as sent", "error", err, "messageID", msg.ID)
			return err
		}
		return nil
	}); err != nil {
		w.logger.Errorw("Transaction failed", "error", err)
		return err
	}

	w.logger.Infow("Message produced and marked as sent", "message", string(msg.Message))

	return nil
}
