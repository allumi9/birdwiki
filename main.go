package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting")

	initPostgresConnection()

	http.HandleFunc("/", welcomePageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
