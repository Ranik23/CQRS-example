package domain



type Order struct {
	ID		    int64		`json:"id"`
	UserID      int64		`json:"user_id"`
	Items 		[]OrderItem `json:"items"`
	TotalPrice  int64		`json:"total_price"`
	Status 		string		`json:"status"`
}

func (o *Order) GetItems() []OrderItem {
	return o.Items
}

func (o *Order) AddItem(item OrderItem) {
	o.Items = append(o.Items, item)
	o.TotalPrice += item.GetPrice()
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

func (o *Order) GetTotalPrice() int64 {
	return o.TotalPrice
}