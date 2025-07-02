package mapper

import (
	"order-service/internal/domain"
	"order-service/internal/domain/dto"
)



func FromRequestCreateOrderToDomainOrder(req dto.RequestCreateOrder) domain.Order {
	order := domain.Order{
		UserID: req.UserID,
		Items:  make([]domain.OrderItem, len(req.Items)),
	}

	for i, item := range req.Items {
		order.Items[i] = domain.OrderItem{
			Name:  item.Name,
			Price: item.Price,
		}
	}

	return order
}