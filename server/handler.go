package server

import (
	"fmt"
	"lamoda_test_task/pkg/models"
	"lamoda_test_task/pkg/repository"
	"log"
	"net/http"
	"strconv"
)

type DefaultArgs struct {
	Message string
}
type CreateArgs struct {
	Products  []models.Product
	Warehosue []models.Warehouse
}

type ProductArgs struct {
	Code    int
	Name    string
	Size    float64
	Value   int64
	StockId string
	Dynamic bool
	ResId   string
}

type WarehouseArgs struct {
	Id        string
	Nmae      string
	Available bool
}

type Response struct {
	Message string
}

func (s *Service) CreateNewProduct(r *http.Request, args *CreateArgs, response *Response) error {
	newProduct := make([]models.Product, 0, len(args.Products))
	productRepo := repository.NewRepository(s.db)
	for _, gd := range args.Products {
		id, err := productRepo.CreateNewProduct(gd.Name, int(gd.Size), int(gd.Value))
		if err != nil {
			return fmt.Errorf("error when creating a new good entity in db: %s", err)
		}
		log.Println("ID:", id)
		newProduct = append(newProduct, gd)
	}

	var resp string = "the following entities have been created: "
	for _, v := range newProduct {
		resp = resp + v.Name + " "
	}

	response.Message = resp
	return nil
}

func (s *Service) CreateNewWarehouse(r *http.Request, args *CreateArgs, response *Response) error {
	newWarehouse := make([]models.Warehouse, 0, len(args.Warehosue))
	productRepo := repository.NewRepository(s.db)
	for _, st := range args.Warehosue {
		if _, err := productRepo.CreateNewWarehouse(st); err != nil {
			return fmt.Errorf("error when creating a new stock entity in db: %s", err)
		}
		newWarehouse = append(newWarehouse, st)
	}

	var resp string = "the following entities have been created: "
	for _, v := range newWarehouse {
		resp = resp + v.Name + " "
	}

	response.Message = resp
	return nil
}

func (s *Service) AddProduct(r *http.Request, args *ProductArgs, response *Response) error {
	productRepo := repository.NewRepository(s.db)
	if err := productRepo.AddProduct(args.Code, args.StockId, args.Value, args.Dynamic); err != nil {
		return fmt.Errorf("error when adding: %s", err)
	}

	response.Message = "done"
	return nil
}

func (s *Service) ReservationProduct(r *http.Request, args *ProductArgs, response *Response) error {
	productRepo := repository.NewRepository(s.db)
	if err := productRepo.ReservationProduct(args.Code, args.StockId, args.Value); err != nil {
		return fmt.Errorf("error when reservation: %s", err)
	}

	response.Message = "done"
	return nil
}

func (s *Service) CancelProductReservation(r *http.Request, args *ProductArgs, response *Response) error {
	productRepo := repository.NewRepository(s.db)
	if err := productRepo.CancelProductReservation(args.ResId); err != nil {
		return fmt.Errorf("error when cancel reservation: %s", err)
	}

	response.Message = "done"
	return nil
}

func (s *Service) GetAllProducts(r *http.Request, args *DefaultArgs, response *Response) error {
	productRepo := repository.NewRepository(s.db)
	goods, err := productRepo.GetProductList()
	if err != nil {
		return fmt.Errorf("error when getting all goods: %s", err)
	}

	var resp string = "all goods: "
	for _, v := range goods {
		resp = resp + v.Name
	}

	response.Message = resp

	return nil
}

func (s *Service) GetProductByCode(r *http.Request, args *ProductArgs, response *Response) error {
	productRepo := repository.NewRepository(s.db)
	g, err := productRepo.GetProduct(args.Code)
	if err != nil {
		return fmt.Errorf("error when getting good by id: %s", err)
	}

	response.Message = ("code: " + string(rune(g.Code)) + " name: " + g.Name + " size: " + strconv.FormatFloat(g.Size, 'f', 2, 64) + " value: " + string(rune(g.Value)))

	return nil
}

func (s *Service) GetProductssCountByWarehouseId(r *http.Request, args *ProductArgs, response *Response) error {
	productRepo := repository.NewRepository(s.db)
	count, err := productRepo.GetProductsCountByWarehouseId(args.StockId, args.Code)
	if err != nil {
		return fmt.Errorf("error when getting goods count by stock id: %s", err)
	}

	response.Message = string(rune(count))

	return nil
}
