package domain

import "errors"

type Order struct {
	ID         int64       `json:"id"`
	UserID     int64       `json:"user_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice int64       `json:"total_price"`

	events []Event `json:"-"`
}

func NewOrder(events []Event) *Order {
	order := &Order{}

	for _, event := range events {
		order.On(event)
	}

	return order
}


func (o *Order) On(event Event) {
	o.events = append(o.events, event)
	switch e := event.(type) {
	case OrderCreatedEvent:
		o.ID = e.OrderID
		o.UserID = e.CustomerID
		o.Items = e.Items
		o.TotalPrice = e.TotalAmount
	case OrderDeletedEvent:
		o = nil // warning
	case ItemAddedEvent:
		o.AddItem(e.Item)
	}
}

func (o *Order) GetItems() []OrderItem {
	return o.Items
}

func (o *Order) GetTotalPrice() int64 {
	return o.TotalPrice
}

func (o *Order) GetEvents() []Event {
	return o.events
}

func (o *Order) AddItem(item OrderItem) error {
	if item.GetPrice() <= 0 {
		return errors.New("item price must be greater than zero")
	}

	if item.GetName() == "" {
		return errors.New("item name cannot be empty")
	}

	o.Items = append(o.Items, item)
	o.TotalPrice += item.GetPrice()

	return nil
}

func (o *Order) RemoveItem(item OrderItem) {
	for i, existingItem := range o.Items {
		if existingItem.GetName() == item.GetName() {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.TotalPrice -= existingItem.GetPrice()
			break
		}
	}
}

func (o *Order) CalculateTotalPrice() {
	o.TotalPrice = 0
	for _, item := range o.Items {
		o.TotalPrice += item.GetPrice()
	}
}
