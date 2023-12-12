package services

import (
	"errors"
	"lamoda_test_task/pkg/models"
	"lamoda_test_task/pkg/repository"
)

type WarehosueService struct {
	repo *repository.Repository
}

func NewWarehouseService(repo *repository.Repository) *WarehosueService {
	return &WarehosueService{repo: repo}
}

// CreateNewWarehouse(name string, available bool) (int, error)
// GetWarehouse(id int) (models.Warehouse, error)
// GetWarehouseList() ([]models.Warehouse, error)

func (w *WarehosueService) CreateNewWarehouse(name string, available bool) (int, error) {
	if name == "" {
		return 0, errors.New("Warehouse name have not value")
	}
	id, err := w.repo.WarehouseRepository.CreateNewWarehouse(name, available)

	if err != nil {
		return 0, nil
	}

	return id, nil
}

func (w *WarehosueService) GetWarehouseByID(id int) (models.Warehouse, error) {
	var wh models.Warehouse
	wh, err := w.repo.WarehouseRepository.GetWarehouse(id)
	if err != nil {
		return models.Warehouse{}, err
	}

	return wh, nil
}

func (w *WarehosueService) GetWarehouseList() ([]models.Warehouse, error) {
	var whs []models.Warehouse
	whs, err := w.repo.WarehouseRepository.GetWarehouseList()

	if err != nil {
		return []models.Warehouse{}, err
	}

	return whs, nil
}
