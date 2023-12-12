package repository

import (
	"database/sql"
	"lamoda_test_task/pkg/models"
	"lamoda_test_task/pkg/repository/postgres"
)

type WarehouseRepository interface {
	CreateNewWarehouse(name string, available bool) (int, error)
	GetWarehouse(id int) (models.Warehouse, error)
	GetWarehouseList() ([]models.Warehouse, error)
}

type ProductRepository interface {
	CreateNewProduct(name string, size float64, value int, stock_id int) (int, error)
	GetProductsCountByWarehouseId(stockId int, code int) (int, error)

	GetProductList() ([]models.Product, error)
	GetProduct(code int) (models.Product, error)

	ReservationProduct(code int, stockId int, value int) error
	CancelProductReservation(resId int) (int, error)
}
type Repository struct {
	WarehouseRepository
	ProductRepository
}

func NewRepository(client *sql.DB) *Repository {
	return &Repository{
		ProductRepository:   postgres.NewProductPSQL(client),
		WarehouseRepository: postgres.NewWarehousePSQL(client),
	}
}
