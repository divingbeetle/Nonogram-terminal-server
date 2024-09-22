package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/divingbeetle/Nonogram-terminal-server/api"
	"github.com/divingbeetle/Nonogram-terminal-server/storage"
)

func main() {
	listenAddr := flag.String("listen-addr", ":8080", "server listen address")

	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASSWORD")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	flag.Parse()

	server := api.NewServer(*listenAddr)
	fmt.Printf("Starting server on %v\n", *listenAddr)

	err := storage.InitDB(dbuser, dbpass, dbhost, dbport, dbname)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())
}
