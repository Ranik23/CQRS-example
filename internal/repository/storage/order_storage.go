package storage

import (
	"context"
	"order-service/internal/domain"
)


//go:generate mockgen -destination=./mock/order_storage.go -package=mock order-service/internal/repository/storage OrderStorage
type OrderStorage interface {
	SaveOrder(ctx context.Context, order domain.Order) error
	GetOrderByID(ctx context.Context, orderID int64) (*domain.Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]domain.Order, error)
	DeleteOrder(ctx context.Context, orderID int64) error
}