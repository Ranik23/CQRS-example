package usecase

import (
	"context"
	"encoding/json"
	"order-service/internal/domain"
	"order-service/internal/repository/storage"
	"order-service/internal/usecase/errs"
	"order-service/pkg/logger"
	"order-service/pkg/txmanager"
)


type DeleteOrderUseCase interface {
	Execute(ctx context.Context, orderID int64) error
}


type deleteOrderUseCase struct {
	mainOrderStorage storage.OrderStorage
	outboxStorage    storage.OutboxStorage
	txmanager        txmanager.TxManager
	logger           logger.Logger
}


func NewDeleteOrderUseCase(mainOrderStorage storage.OrderStorage, outboxStorage storage.OutboxStorage, txmanager txmanager.TxManager, logger logger.Logger) *deleteOrderUseCase {
	return &deleteOrderUseCase{
		mainOrderStorage: mainOrderStorage,
		outboxStorage:    outboxStorage,
		txmanager:        txmanager,
		logger:           logger,
	}
}


func (u *deleteOrderUseCase) Execute(ctx context.Context, orderID int64) error {
	if orderID <= 0 {
		u.logger.Errorw("Invalid order ID", "order_id", orderID)
		return errs.ErrInvalidOrderID
	}

	if err := u.txmanager.Run(ctx, func(ctx context.Context) error {

		order, err := u.mainOrderStorage.GetOrderByID(ctx, orderID)
		if err != nil {
			u.logger.Errorw("Failed to get order by ID", "error", err, "order_id", orderID)
			return err
		}

		if err := u.mainOrderStorage.DeleteOrder(ctx, orderID); err != nil {
			u.logger.Errorw("Failed to delete order", "error", err, "order_id", orderID)
			return err
		}

		orderBytes, err := json.Marshal(order)
		if err != nil {
			u.logger.Errorw("Failed to marshal order", "error", err, "order_id", orderID)
			return err
		}

		if err := u.outboxStorage.CreateOutboxMessage(ctx, domain.OrderDeletedKey, orderBytes); err != nil {
			u.logger.Errorw("Failed to create outbox message for deleted order", "error", err, "order_id", orderID)
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorw("Transaction failed", "error", err, "order_id", orderID)
		return err
	}

	return nil
}