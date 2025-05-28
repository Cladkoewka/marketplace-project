package event

type OrderPlacedEvent struct {
	Event      string `json:"event"`
	Timestamp  string `json:"timestamp"`
	OrderID    int64  `json:"order_id"`
	CustomerID int64  `json:"customer_id"`
	Status     string `json:"status"`
}
