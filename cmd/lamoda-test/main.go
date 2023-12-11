package main

import (
	"lamoda_test_task/config"
	"lamoda_test_task/pkg/repository/postgres"
	"lamoda_test_task/server"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	log.Println("Starting rpc server")
	config.Init()
	log.Println("Connecting to database")
	conn, err := postgres.EstablishPSQLConnection(&postgres.PSQLConfig{
		Host:     viper.GetString("db.postgres.host"),
		Port:     viper.GetString("db.postgres.port"),
		Password: viper.GetString("db.postgres.password"),
		DBName:   viper.GetString("db.postgres.database"),
		Username: viper.GetString("db.postgres.user"),
		SSLMode:  viper.GetString("db.postgres.sslmode"),
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connection establish")

	log.Printf("Server listing on port: %s", viper.GetString("app.port"))
	server.RunJRPC(conn, viper.GetString("app.port"))
}
