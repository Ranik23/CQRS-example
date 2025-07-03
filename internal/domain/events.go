package domain

import "time"

type Event interface {
	isEvent()
}

type OrderCreatedEvent struct {
	OrderID 		int64 `json:"order_id"`
	CustomerID 		int64 `json:"customer_id"`
	Items 			[]OrderItem `json:"items"`
	TotalAmount 	int64 `json:"total_amount"`
}

type OrderDeletedEvent struct {
	OrderID 		int64 `json:"order_id"`
	DeletedAt 		time.Time `json:"deleted_at"`
}

type ItemAddedEvent struct {
	OrderID int64 `json:"order_id"`
	Item 	OrderItem `json:"item"`
	AddedAt time.Time `json:"added_at"`
}



func (ItemAddedEvent) isEvent() {}

func (OrderCreatedEvent) isEvent() {}

func (OrderDeletedEvent) isEvent() {}