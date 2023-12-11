package server

import (
	"database/sql"
	"lamoda_test_task/pkg/repository"
	"lamoda_test_task/pkg/services"
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type Service struct {
	service *services.Service
	// ctx context.Context
}

func NewService(service *services.Service) *Service {
	return &Service{
		service: service,
	}
}

func RunJRPC(db *sql.DB, port string) {
	s := rpc.NewServer()
	log.Println("run server")
	// ctx := context.Background()
	// services := NewService(delivery.NewHandler(services.NewService(repository.NewRepository(db))))
	// ProductService := delivery.NewHandler(services.NewService(repository.NewRepository(services.db)))
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(NewService(services.NewService(repository.NewRepository(db))), "")
	log.Println("register service")

	http.Handle("/rpc", s)

	http.ListenAndServe(":"+port, s)
	log.Printf("Listen and Serve on %s", port)
}
