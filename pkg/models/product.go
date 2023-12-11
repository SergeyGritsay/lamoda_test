package models

type Product struct {
	Code    int64   `json:"good_code"`
	Name    string  `json:"good_name"`
	Size    float64 `json:"good_size"`
	Value   int64   `json:"good_value"`
	StockId int     `json:"warehouse_id"`
}
