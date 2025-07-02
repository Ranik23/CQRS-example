package consumer

import "context"



//go:generate mockgen -destination=./mock/consumer.go -package=mock order-service/internal/infrastructure/consumer Consumer
type Consumer interface {
	Consume(ctx context.Context) ([]byte, []byte, error)
}


