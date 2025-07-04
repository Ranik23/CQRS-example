package worker

import (
	"context"
	"errors"
	"order-service/internal/infrastructure/consumer"
	"order-service/internal/infrastructure/producer"
	"order-service/internal/models"
	"order-service/internal/repository/storage"
	"order-service/pkg/logger"
	"order-service/pkg/txmanager"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

type Worker struct {
	outboxStorage storage.OutboxStorage
	producer      producer.Producer
	consumer      consumer.Consumer
	txmanager     txmanager.TxManager
	pool          *pgxpool.Pool

	logger logger.Logger
	ctx    context.Context
	cancel context.CancelFunc

	numWorker int
	batchSize int
}

func NewWorker(producer producer.Producer, consumer consumer.Consumer, outboxStorage storage.OutboxStorage,
	txmanager txmanager.TxManager, logger logger.Logger, pool *pgxpool.Pool, numWorker int, batchSize int) *Worker {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Worker{
		producer:      producer,
		consumer:      consumer,
		outboxStorage: outboxStorage,
		txmanager:     txmanager,
		logger:        logger,
		ctx:           ctx,
		cancel:        cancelFunc,
		pool:          pool,
		numWorker:     numWorker,
		batchSize:     batchSize,
	}
}

func (w *Worker) Run() error {
	w.logger.Infow("Worker started")
	defer w.logger.Infow("Worker stopped")

	g, errctx := errgroup.WithContext(context.Background())

	for workerID := 0; workerID < w.numWorker; workerID++ {
		g.Go(func() error {
			w.logger.Infow("Starting dispatch worker", "workerID", workerID)
			if err := w.dispatchWorker(w.ctx, errctx, workerID); err != nil && !errors.Is(err, context.Canceled) {
				w.logger.Errorw("Dispatch Worker failed", "error", err, "workerID", workerID)
				return err
			}
			w.logger.Infow("Dispatch Worker finished", "workerID", workerID)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		w.logger.Errorw("Worker encountered an error", "error", err)
		return err
	}

	return nil
}

func (w *Worker) Stop() {
	w.cancel()
	w.logger.Infow("Worker stopped")
}

func (w *Worker) dispatchWorker(ctx context.Context, errctx context.Context, workerID int) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	w.logger.Infow("Dispatch Worker started", "workerID", workerID)
	for {
		select {
		case <-ctx.Done():
			w.logger.Infow("Dispatch Worker stopping due to ctx done", "workerID", workerID)
			return context.Canceled
		case <-errctx.Done():
			w.logger.Infow("Dispatch Worker stopping due to errgroup ctx done", "workerID", workerID)
			return errctx.Err()
		case <-ticker.C:
			w.logger.Debugw("Worker tick", "workerID", workerID)
			if err := w.dispatchEvent(ctx); err != nil {
				w.logger.Errorw("Dispatch event failed", "error", err, "workerID", workerID)
			} else {
				w.logger.Debugw("Dispatch event succeeded", "workerID", workerID)
			}
		}
	}
}
func (w *Worker) dispatchEvent(ctx context.Context) error {
	w.logger.Debug("Attempting to dispatch event")

	tx, err := w.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		w.logger.Errorw("Failed to begin transaction", "error", err)
		return err
	}
	defer func() {
		if tx != nil {
			if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				w.logger.Errorw("Failed to rollback transaction", "error", err)
			}
		}
	}()

	rows, err := tx.Query(ctx,
		`
		WITH next_messages AS (
			SELECT id
			FROM outbox
			WHERE status IN ('not sent')
			ORDER BY created_at ASC
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE outbox
		SET status = 'processing'
		WHERE id IN (SELECT id FROM next_messages)
		RETURNING id, status, key, message
		`, w.batchSize)

	if err != nil {
		return err
	}
	defer rows.Close()

	var msgs []models.OutboxMessage
	for rows.Next() {
		var m models.OutboxMessage
		if err := rows.Scan(&m.ID, &m.Sent, &m.Key, &m.Message); err != nil {
			return err
		}
		msgs = append(msgs, m)
	}

	if len(msgs) == 0 {
		w.logger.Debugw("No outbox messages to process")
		tx = nil
		return nil
	}

	if err := tx.Commit(ctx); err != nil {
		w.logger.Errorw("Failed to commit transaction", "error", err)
		return err
	}
	tx = nil

	// Отправляем сообщения в брокер
	for _, msg := range msgs {
		if err := w.producer.Produce(ctx, []byte(msg.Key), msg.Message); err != nil {
			w.logger.Errorw("Failed to produce message to broker", "error", err, "messageID", msg.ID)
			return err
		}
		w.logger.Infow("Produced message to broker", "messageID", msg.ID)
	}

	// Обновляем статус всех сообщений сразу
	tx2, err := w.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		w.logger.Errorw("Failed to begin transaction for marking as sent", "error", err)
		return err
	}
	defer func() {
		if tx2 != nil {
			if err := tx2.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
				w.logger.Errorw("Failed to rollback transaction", "error", err)
			}
		}
	}()

	// Собираем все ID сообщений в срез
	ids := make([]int64, 0, len(msgs))
	for _, m := range msgs {
		ids = append(ids, m.ID)
	}

	// Обновляем статусы батчем
	_, err = tx2.Exec(ctx,
		`UPDATE outbox SET status = 'sent' WHERE id = ANY($1)`,
		ids)
	if err != nil {
		w.logger.Errorw("Failed to update outbox messages status to sent", "error", err)
		return err
	}

	if err := tx2.Commit(ctx); err != nil {
		w.logger.Errorw("Failed to commit transaction for marking as sent", "error", err)
		return err
	}
	tx2 = nil

	w.logger.Infow("Marked outbox messages as sent", "count", len(msgs))

	return nil
}
