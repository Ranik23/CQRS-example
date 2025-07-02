package domain


type OrderItem struct {
	ID    	int64	`json:"id"`
	Name  	string  `json:"name"`
	Price 	int64   `json:"price"`
}

func (o *OrderItem) GetName() string {
	return o.Name
}

func (o *OrderItem) GetPrice() int64 {
	return o.Price
}