package usecase

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/repository/storage"
	"order-service/internal/usecase/errs"
)


//go:generate mockgen -destination=./mock/get_orders_usecase.go -package=mock order-service/internal/usecase GetOrdersUseCase
type GetOrdersUseCase interface {
	Execute(ctx context.Context, userID int64) ([]domain.Order, error)
}

type getOrdersUseCase struct {
	mainOrderStorage storage.OrderStorage
}

func NewGetOrdersUseCase(mainOrderStorage storage.OrderStorage) *getOrdersUseCase {
	return &getOrdersUseCase{mainOrderStorage: mainOrderStorage}
}

func (u *getOrdersUseCase) Execute(ctx context.Context, userID int64) ([]domain.Order, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidUserID
	}

	orders, err := u.mainOrderStorage.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}