package models

type Warehouse struct {
	ID          int    `json:"stock_id"`
	Name        string `json:"stock_name"`
	IsAvailable bool   `json:"stock_available"`
}
