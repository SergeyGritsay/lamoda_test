package repository

import (
	"database/sql"
	"lamoda_test_task/pkg/models"
	"lamoda_test_task/pkg/repository/postgres"
)

type WarehosueRepository interface {
	CreateNewWarehouse(s models.Warehouse) (int, error)
	GetWarehouse(id int) (models.Warehouse, error)
	GetWarehouseList() ([]models.Warehouse, error)
}

type ProductRepository interface {
	CreateNewProduct(name string, size int, value int) (int, error)
	GetProductsCountByWarehouseId(stockId string, code int) (int64, error)

	GetProductList() ([]models.Product, error)
	GetProduct(code int) (models.Product, error)

	ReservationProduct(code int, stockId string, value int64) error
	CancelProductReservation(resId string) error
	AddProduct(code int, stockId string, value int64, dynamic bool) error
}
type Repository struct {
	WarehosueRepository
	ProductRepository
}

func NewRepository(client *sql.DB) *Repository {
	return &Repository{
		ProductRepository:   postgres.NewProductPSQL(client),
		WarehosueRepository: postgres.NewWarehousePSQL(client),
	}
}
