package handlers

import (
	"fmt"
	"order-service/internal/domain/dto"
	"order-service/internal/mapper"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with user ID and items
// @Tags orders
// @Accept json
// @Produce json
// @Param order body RequestCreateOrder true "Order data"
// @Success 201 {object} map[string]string "Order created successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Failed to create order"
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	fmt.Println("Handler CreateOrder called")
	var req dto.RequestCreateOrder
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		h.Logger.Errorw("Failed to bind request", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if err := h.CreateOrderUseCase.Execute(c, mapper.FromRequestCreateOrderToDomainOrder(req)); err != nil {
		h.Logger.Errorw("Failed to create order", "error", err)
		c.JSON(500, gin.H{"error": "Failed to create order"})
		return 
	}

	h.Logger.Infow("Order created successfully", "user_id", req.UserID, "items", req.Items)

	c.JSON(200, gin.H{"message": "Order created successfully"})
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an existing order by its ID
// @Tags orders
// @Param id path int true "Order ID"
// @Produce json
// @Success 200 {object} map[string]string "Order deleted successfully"
// @Failure 400 {object} map[string]string "Invalid Order ID"
// @Failure 500 {object} map[string]string "Failed to delete order"
// @Router /orders/{id} [delete]
func (h *Handler) DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		h.Logger.Errorw("Order ID is required")
		c.JSON(400, gin.H{"error": "Order ID is required"})
		return
	}

	order_id, err := strconv.Atoi(orderID)
	if err != nil {
		h.Logger.Errorw("Invalid Order ID", "error", err)
		c.JSON(400, gin.H{"error": "Invalid Order ID"})
		return
	}

	if err := h.DeleteOrderUseCase.Execute(c, int64(order_id)); err != nil {
		h.Logger.Errorw("Failed to delete order", "error", err, "order_id", orderID)
		c.JSON(500, gin.H{"error": "Failed to delete order"})
		return
	}
	h.Logger.Infow("Order deleted successfully", "order_id", orderID)
	c.JSON(200, gin.H{"message": "Order deleted successfully"})
}