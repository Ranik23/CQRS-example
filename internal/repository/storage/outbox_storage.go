package storage

import (
	"context"
	"order-service/internal/models"
)

//go:generate mockgen -destination=./mock/outbox_storage.go -package=mock order-service/internal/repository/storage OutboxStorage
type OutboxStorage interface {
	CreateOutboxMessage(ctx context.Context, key string, message []byte) error
	GetOutBoxMessage(ctx context.Context) (*models.OutboxMessage, error)
	MarkAsSent(ctx context.Context, id int64) error
}

