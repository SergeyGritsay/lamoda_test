package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PSQLConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

const driverName = "postgres"

func EstablishPSQLConnection(cnf *PSQLConfig) (*sql.DB, error) {
	log.Println("Starting connection to db")
	time.Sleep(10 * time.Second)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cnf.Host, cnf.Port, cnf.Username, cnf.Password, cnf.DBName, cnf.SSLMode)
	db, err := sql.Open(driverName, psqlInfo)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(fmt.Sprintf("Connected to db: %s", cnf.DBName))

	return db, nil
}
