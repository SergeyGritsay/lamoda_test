package server

import (
	"context"
	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

type Service struct {
	ctx context.Context
	// db  db.Storage
	// log *logger.Logger
}

// func NewService(ctx context.Context, strg db.Storage, log *logger.Logger) *Service {
func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
		// db:  strg,
		// log: log,
	}
}

// func RunJRPC(ctx context.Context, strg db.Storage, log *logger.Logger) {
func RunJRPC(ctx context.Context, port string) {
	s := rpc.NewServer()
	// log.Info("run server")

	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(NewService(ctx), "") // Add db and logger
	// log.Info("register service")

	http.Handle("/rpc", s)

	http.ListenAndServe(":"+port, nil)
}
