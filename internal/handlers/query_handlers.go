package handlers

import "github.com/gin-gonic/gin"



func (h *Handler) GetOrderItems(c *gin.Context){
	var req RequestGetOrderItems
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	orderItems, err := h.GetOrdersUseCase.Execute(c, req.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve order items"})
		return
	}

	c.JSON(200, gin.H{"order_items": orderItems})
}