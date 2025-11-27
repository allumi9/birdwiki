package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting")

	http.HandleFunc("/", welcomePage)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
