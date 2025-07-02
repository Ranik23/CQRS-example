package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/internal/domain"
	"order-service/internal/repository/cache"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client redis.Client
}

// GetOrdersByUserID implements cache.Cache.
func (r *redisCache) GetOrdersByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	key := "orders:" + fmt.Sprint(userID)
	orderData, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var orders []domain.Order
	for _, data := range orderData {
		var order domain.Order
		if err := json.Unmarshal([]byte(data), &order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// CreateOrder implements cache.Cache.
func (r *redisCache) CreateOrder(ctx context.Context, order *domain.Order) error {
	orderData, err := json.Marshal(*order)
	if err != nil {
		return err
	}

	key := "orders:" + fmt.Sprint(order.UserID)
	if err := r.client.RPush(ctx, key, orderData).Err(); err != nil {
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
