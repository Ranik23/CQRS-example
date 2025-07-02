package handlers


type CreateOrderItem struct {
	Name  string `json:"name" binding:"required"`
	Price int64  `json:"price" binding:"required,gt=0"`
}

type RequestCreateOrder struct {
	UserID int64             `json:"user_id" binding:"required,gt=0"`
	Items  []CreateOrderItem `json:"items" binding:"required,min=1,dive"`
}

type RequestGetOrderItems struct {
	UserID int64 `json:"customer_id" binding:"required"`
}

