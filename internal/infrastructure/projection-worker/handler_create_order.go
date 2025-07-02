package projectionworker

import (
	"context"
	"encoding/json"
	"order-service/internal/domain"
)



func (w *ProjectionWorker) handleCreateOrder(ctx context.Context, value []byte) error {
	var order domain.Order
	if err := json.Unmarshal(value, &order); err != nil {
		w.logger.Errorw("Failed to unmarshal order", "error", err)
		return err
	}

	if err := w.cacheClient.CreateOrder(ctx, &order); err != nil {
		w.logger.Errorw("Failed to create order in cache", "error", err, "order_id", order.ID)
		return err
	}

	return nil
}