package delivery

import (
	"encoding/json"
	"fmt"
	"lamoda_test_task/pkg/models"
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
	StockId int
	Dynamic bool
	ResId   int
}

type ReservationArgs struct {
	Codes   []int
	Value   []int
	StockId int
}

type WarehouseArgs struct {
	Id        string
	Name      string
	Available bool
}

type Response struct {
	Message string
}

func (h *Handler) CreateNewProduct(r *http.Request, args *ProductArgs, response *Response) error {

	id, err := h.services.Product.CreateNewProduct(args.Name, args.Size, int(args.Value), args.StockId)
	if err != nil {
		return fmt.Errorf("error when creating a new good entity in db: %s", err)
	}

	var resp string = "the following entities have been created: " + strconv.Itoa(id) + " " + args.Name

	response.Message = resp
	return nil
}

func (h *Handler) CreateNewWarehouse(r *http.Request, args *WarehouseArgs, response *Response) error {
	// productRepo := repository.NewRepository(s.db)
	id, err := h.services.Warehouse.CreateNewWarehouse(args.Name, args.Available)
	if err != nil {
		return fmt.Errorf("error when creating a new stock entity in db: %s", err)
	}

	var resp string = "the following entities have been created: "
	resp = resp + " " + strconv.Itoa(id) + args.Name + " "

	response.Message = resp
	return nil
}

func (h *Handler) ReservationProduct(r *http.Request, args *ReservationArgs, response *Response) error {
	if err := h.services.Product.ReservationProducts(args.Codes, args.StockId, args.Value); err != nil {
		return fmt.Errorf("error when reservation: %s", err)
	}

	response.Message = "done"
	return nil
}

func (h *Handler) CancelProductReservation(r *http.Request, args *ProductArgs, response *Response) error {
	code, err := h.services.Product.CancelProductReservation(args.ResId)
	if err != nil {
		return fmt.Errorf("error when cancel reservation: %s", err)
	}

	response.Message = "done. Code: " + strconv.Itoa(code)
	return nil
}

func (h *Handler) GetAllProducts(r *http.Request, args *DefaultArgs, response *Response) error {
	goods, err := h.services.Product.GetProductList()
	if err != nil {
		return fmt.Errorf("error when getting all goods: %s", err)
	}

	var resp string = "all goods: "
	for _, v := range goods {
		resp = resp + " " + strconv.Itoa(int(v.Code)) + " " + v.Name + " " + fmt.Sprintf("%.6f", v.Size) + " " + strconv.Itoa(int(v.Value)) + " " + strconv.Itoa(v.StockId) + "\n"
	}

	response.Message = resp

	return nil
}

func (h *Handler) GetProductByCode(r *http.Request, args *ProductArgs, response *Response) error {
	product, err := h.services.Product.GetProductByCode(args.Code)
	if err != nil {
		return fmt.Errorf("error when getting good by id: %s", err)
	}
	js, err := json.Marshal(product)
	response.Message = string(js)

	return nil
}

func (h *Handler) GetProductssCountByWarehouseId(r *http.Request, args *ProductArgs, response *Response) error {
	count, err := h.services.Product.GetProductsCountByWarehouseId(args.StockId, args.Code)
	if err != nil {
		return fmt.Errorf("error when getting goods count by stock id: %s", err)
	}

	response.Message = string(rune(count))

	return nil
}
