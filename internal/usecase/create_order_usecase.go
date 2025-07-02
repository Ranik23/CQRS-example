package usecase

import (
	"context"
	"encoding/json"

	"order-service/internal/domain"
	"order-service/internal/repository/storage"
	"order-service/pkg/logger"
	"order-service/pkg/txmanager"
)

//go:generate mockgen -destination=./mock/create_order_usecase.go -package=mock order-service/internal/usecase CreateOrderUseCase
type CreateOrderUseCase interface {
	Execute(ctx context.Context, order domain.Order) error
}

type createOrderUseCase struct {
	mainOrderStorage storage.OrderStorage
	outboxStorage    storage.OutboxStorage
	txmanager        txmanager.TxManager
	logger           logger.Logger
}

func NewCreateOrderUseCase(mainOrderStorage storage.OrderStorage, outboxStorage storage.OutboxStorage, txmanager txmanager.TxManager, logger logger.Logger) *createOrderUseCase {
	return &createOrderUseCase{mainOrderStorage: mainOrderStorage, outboxStorage: outboxStorage, txmanager: txmanager, logger: logger}
}

func (u *createOrderUseCase) Execute(ctx context.Context, order domain.Order) error {
	if err := u.txmanager.Run(ctx, func(ctx context.Context) error {
		if err := u.mainOrderStorage.SaveOrder(ctx, order); err != nil {
			u.logger.Errorw("Failed to save order", "error", err)
			return err
		}
		u.logger.Infow("Order saved", "order_id", order.ID, "user_id", order.UserID)

		orderBytes, err := json.Marshal(order)
		if err != nil {
			u.logger.Errorw("Failed to marshal order", "error", err)
			return err
		}

		if err := u.outboxStorage.CreateOutboxMessage(ctx, domain.OrderCreatedKey, orderBytes); err != nil {
			u.logger.Errorw("Failed to create outbox message", "error", err)
			return err
		}
		
		u.logger.Infow("Outbox message created", "order_id", order.ID)

		return nil

	}); err != nil {
		u.logger.Errorw("Transaction failed", "error", err)
		return err
	}

	u.logger.Infow("Order created successfully", "order_id", order.ID, "user_id", order.UserID)

	return nil
}
