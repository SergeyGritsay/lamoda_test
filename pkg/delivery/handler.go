package delivery

import "lamoda_test_task/pkg/services"

type Handler struct {
	services *services.Service
}

func NewHandler(serv *services.Service) *Handler {
	return &Handler{services: serv}
}
