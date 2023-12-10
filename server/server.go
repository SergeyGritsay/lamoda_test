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
	ctx context.Context
	db  *sql.DB
	// log *logger.Logger
}

// func NewService(ctx context.Context, strg db.Storage, log *logger.Logger) *Service {
func NewService(ctx context.Context, db *sql.DB) *Service {
	return &Service{
		ctx: ctx,
		db:  db,
		// log: log,
	}
}

// func RunJRPC(ctx context.Context, strg db.Storage, log *logger.Logger) {
func RunJRPC(ctx context.Context, db *sql.DB, port string) {
	s := rpc.NewServer()
	// log.Info("run server")

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(NewService(ctx, db), "") // Add db and logger
	// log.Info("register service")

	http.Handle("/rpc", s)

	http.ListenAndServe(":"+port, s)
	log.Printf("Listen and Serve on %s", port)
}
