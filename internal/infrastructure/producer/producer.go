package producer

import "context"



//go:generate mockgen -destination=./mock/producer.go -package=mock order-service/internal/infrastructure/producer Producer
type Producer interface {
	Produce(ctx context.Context, key []byte, message []byte) error
}

