package models

type (
	Warehouse struct {
		ID          string `json:"stock_id"`
		Name        string `json:"stock_name"`
		IsAvailable bool   `json:"stock_available"`
	}

	Product struct {
		Code  int64   `json:"good_code"`
		Name  string  `json:"good_name"`
		Size  float64 `json:"good_size"`
		Value int64   `json:"good_value"`
	}
)
