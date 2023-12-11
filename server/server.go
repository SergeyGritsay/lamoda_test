package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type Service struct {
	db  *sql.DB
	ctx context.Context
}

func NewService(db *sql.DB, ctx context.Context) *Service {
	return &Service{
		db:  db,
		ctx: ctx,
	}
}

func RunJRPC(db *sql.DB, port string) {
	s := rpc.NewServer()
	log.Println("run server")
	ctx := context.Background()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(NewService(db, ctx), "")
	log.Println("register service")

	http.Handle("/rpc", s)

	http.ListenAndServe(":"+port, s)
	log.Printf("Listen and Serve on %s", port)
}
