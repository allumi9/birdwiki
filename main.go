package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting")

	defer closePostgresConnection(postgresConnection)

	http.HandleFunc("/", welcomePageHandler)
	http.HandleFunc("/api/search", searchRequestHandler)
	http.HandleFunc("/bird-info", birdInfoPageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
