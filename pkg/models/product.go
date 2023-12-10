package models

type Product struct {
	Id              uint   `json:"product_id"`       // Unique code
	ProductName     string `json:"product_name"`     // Name
	ProductSize     string `json:"product_size"`     // Size
	ProductQuantity uint   `json:"product_quantity"` // Amount of products
}
