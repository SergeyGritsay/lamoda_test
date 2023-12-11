package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func RunJRPC(db *sql.DB, port string) {
	s := rpc.NewServer()
	log.Println("run server")

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(NewService(db), "")
	log.Println("register service")

	http.Handle("/rpc", s)

	http.ListenAndServe(":"+port, s)
	log.Printf("Listen and Serve on %s", port)
}
