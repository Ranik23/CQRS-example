package redis

import (
	"fmt"
	"context"
	"encoding/json"
	"order-service/internal/domain"
	"order-service/internal/repository/cache"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client redis.Client
}

// CreateOrder implements cache.Cache.
func (r *redisCache) CreateOrder(ctx context.Context, order *domain.Order) error {
	orderData, err := json.Marshal(*order)
	if err != nil {
		return err
	}

	key := "order:" + fmt.Sprint(order.ID)
	err = r.client.Set(ctx, key, orderData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// DeleteOrder implements cache.Cache.
func (r *redisCache) DeleteOrder(ctx context.Context, orderID int64) error {
	key := "order:" + fmt.Sprint(orderID)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetOrder implements cache.Cache.
func (r *redisCache) GetOrder(ctx context.Context, orderID int64) (*domain.Order, error) {
	key := "order:" + fmt.Sprint(orderID)
	orderData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var order domain.Order
	err = json.Unmarshal([]byte(orderData), &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func NewRedisCache(address string) cache.Cache {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	return &redisCache{client: *client}
}
