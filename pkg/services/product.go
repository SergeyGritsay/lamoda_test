package services

import (
	"errors"
	"lamoda_test_task/pkg/models"
	"lamoda_test_task/pkg/repository"
)

type ProductService struct {
	repo *repository.Repository
}

func NewProductService(repo *repository.Repository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (p *ProductService) CreateNewProduct(name string, size float64, value int, stockId int) (int, error) {
	if name == "" || size == 0 || value == 0 || stockId == 0 {
		return 0, errors.New("Invalid fields")
	}

	productId, err := p.repo.ProductRepository.CreateNewProduct(name, size, value, stockId)
	if err != nil {
		return 0, nil
	}

	return productId, nil
}

func (p *ProductService) GetProductByCode(code int) (models.Product, error) {
	product, err := p.repo.ProductRepository.GetProduct(code)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (p *ProductService) GetProductList() ([]models.Product, error) {
	var productList []models.Product

	productList, err := p.repo.ProductRepository.GetProductList()

	if err != nil {
		return []models.Product{}, err
	}

	return productList, err
}

func (p *ProductService) GetProductsCountByWarehouseId(stockId int, code int) (int64, error) {
	count, err := p.repo.ProductRepository.GetProductsCountByWarehouseId(stockId, code)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *ProductService) ReservationProducts(code []int, stockId int, value []int64) error {
	if len(code) < len(value) || len(value) < len(code) {
		return errors.New("Missing products or value of products")
	}
	ok := make([]error, 0, len(code))

	for idx, val := range code {
		err := p.repo.ProductRepository.ReservationProduct(val, stockId, value[idx])
		if err != nil {
			ok = append(ok, err)
		}
	}

	if len(ok) != 0 {
		return errors.New("Something wrong in reservaton product")
	}

	return nil
}

func (p *ProductService) CancelProductReservation(resId string) error {
	err := p.repo.ProductRepository.CancelProductReservation(resId)

	if err != nil {
		return err
	}

	return nil
}
