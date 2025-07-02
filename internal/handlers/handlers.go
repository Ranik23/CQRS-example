package handlers

import (
	"order-service/internal/usecase"
	"order-service/pkg/logger"
)



type Handler struct {
	usecase.CreateOrderUseCase
	usecase.GetOrdersUseCase
	usecase.DeleteOrderUseCase
	logger.Logger
}



func NewHandler(createOrderUseCase usecase.CreateOrderUseCase, deleteOrderUseCase usecase.DeleteOrderUseCase, getOrdersUseCase usecase.GetOrdersUseCase, logger logger.Logger) *Handler {
	return &Handler{
		CreateOrderUseCase: createOrderUseCase,
		DeleteOrderUseCase: deleteOrderUseCase,
		GetOrdersUseCase:   getOrdersUseCase,
		Logger:             logger,
	}
}


