package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flag.String("dir", ".", "directory with migration files")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	args := flag.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	dbstring, command := args[1], args[2]

	db, err := sql.Open("postgres", dbstring)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: Failed to close db: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
