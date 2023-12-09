package main

import (
	"context"
	"lamoda_test_task/server"
	"log"
)

func main() {
	log.Println("Starting rpc server")
	ctx := context.Background()
	log.Println("Server listing on port :3000")

	server.RunJRPC(ctx, "3000")
}
