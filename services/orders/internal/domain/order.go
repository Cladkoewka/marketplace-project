package domain

type Order struct {
	ID         int64    `json:"id"`
	CustomerID int64    `json:"customer_id"`
	Status     string `json:"status"`
}
