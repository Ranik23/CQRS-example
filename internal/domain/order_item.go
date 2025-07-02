package domain


type OrderItem struct {
	Name  	string  `json:"name"`
	Price 	int64   `json:"price"`
}

func (o *OrderItem) GetName() string {
	return o.Name
}

func (o *OrderItem) GetPrice() int64 {
	return o.Price
}