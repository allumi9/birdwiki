package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

const databaseUrl string = "postgresql://postgres:1234@localhost:5432/postgres"

func initPostgresConnection() {
	connection, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Printf("Error connecting to PSQL: %v\n", err)
		return
	}
	defer connection.Close(context.Background())

	var birdName string
	err = connection.QueryRow(context.Background(), "select name from bird_files;").Scan(birdName)
	if err != nil {
		log.Printf("Error querying db: %v", err)
		return
	}

	fmt.Println(birdName)
}

func main() {
	fmt.Println("Starting")

	initPostgresConnection()

	http.HandleFunc("/", welcomePageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
