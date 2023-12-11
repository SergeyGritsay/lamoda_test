package services

import (
	"lamoda_test_task/pkg/models"
	"lamoda_test_task/pkg/repository"
)

type Product interface {
	CreateNewProduct(name string, size float64, value int, stock_id int) (int, error)
	GetProductsCountByWarehouseId(stockId int, code int) (int64, error)

	GetProductList() ([]models.Product, error)
	GetProductByCode(code int) (models.Product, error)

	ReservationProducts(code []int, stockId int, value []int64) error
	CancelProductReservation(resId string) error
}

type Warehouse interface {
	CreateNewWarehouse(name string, available bool) (int, error)
	GetWarehouseByID(id int) (models.Warehouse, error)
	GetWarehouseList() ([]models.Warehouse, error)
}

type Service struct {
	Product
	Warehouse
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Product:   NewProductService(repo),
		Warehouse: NewWarehouseService(repo),
	}
}
