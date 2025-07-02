package cache

import (
	"context"
	"order-service/internal/domain"
)



type Cache interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrder(ctx context.Context, orderID int64) (*domain.Order, error)
	DeleteOrder(ctx context.Context, orderID int64) error
}

