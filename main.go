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
	http.Handle("/bird_pictures/", http.StripPrefix("/bird_pictures/", http.FileServer(http.Dir("./bird_pictures"))))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
