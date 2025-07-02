package usecase

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/repository/cache"
	"order-service/internal/usecase/errs"
)

//go:generate mockgen -destination=./mock/get_orders_usecase.go -package=mock order-service/internal/usecase GetOrdersUseCase
type GetOrdersUseCase interface {
	Execute(ctx context.Context, userID int64) ([]domain.Order, error)
}

type getOrdersUseCase struct {
	cacheStorage cache.Cache
}

func NewGetOrdersUseCase(cache cache.Cache) *getOrdersUseCase {
	return &getOrdersUseCase{cacheStorage: cache}
}

func (u *getOrdersUseCase) Execute(ctx context.Context, userID int64) ([]domain.Order, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidUserID
	}

	orders, err := u.cacheStorage.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	return orders, nil
}