package postgres

import (
	"context"
	"errors"
	"fmt"
	"order-service/internal/models"
	"order-service/internal/repository/storage"
	"order-service/internal/repository/storage/errs"
	"order-service/pkg/txmanager"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type outboxStorage struct {
	txmanager txmanager.TxManager
	sq        squirrel.StatementBuilderType
}

// MarkAsSent implements storage.OutboxStorage.
func (o *outboxStorage) MarkAsSent(ctx context.Context, id int64) error {
	tx := o.txmanager.Tx(ctx)

	query, args, err := o.sq.Update("outbox").
		Set("status", "sent").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// GetOutBoxMessage implements storage.OutboxStorage.
func (o *outboxStorage) GetOutBoxMessage(ctx context.Context) (*models.OutboxMessage, error) {
	tx := o.txmanager.Tx(ctx)

	query, args, err := o.sq.Select("id, status, key, message").
		From("outbox").
		Where(squirrel.Eq{"status": "not sent"}).
		OrderBy("created_at ASC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(ctx, query, args...)
	var outboxMessage models.OutboxMessage
	if err := row.Scan(&outboxMessage.ID, &outboxMessage.Sent, &outboxMessage.Key, &outboxMessage.Message); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNoFound
		}
		return nil, err
	}

	return &outboxMessage, nil
}

func NewOutboxStorage(txManager txmanager.TxManager) storage.OutboxStorage {
	return &outboxStorage{txmanager: txManager, sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func (o *outboxStorage) CreateOutboxMessage(ctx context.Context, key string, message []byte) error {
	tx := o.txmanager.Tx(ctx)

	fmt.Println("Creating outbox message with key:", key, "and message:", string(message))

	query, args, err := o.sq.Insert("outbox").
		Columns("key", "message").
		Values(key, message).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return err
}
