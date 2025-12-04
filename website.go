package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Bird struct {
	Name        string
	Description string
}

type WelcomePage struct {
	Title string
	Birds []Bird
}

func getBird(birdPath string) (Bird, error) {
	fileDescription, err := os.ReadFile(birdPath)
	if err != nil {
		errorMessage := "Error when reading file" + birdPath
		log.Printf(errorMessage, err)
		return Bird{"", ""}, errors.New(errorMessage)
	}

	slicedDesc := strings.Split(string(fileDescription), "\n")
	return Bird{slicedDesc[0], slicedDesc[1]}, nil
}

func getBirdsArray() []Bird {
	var birds []Bird

	jackdaw, _ := queryBirdTableByName("jackdaw")
	crow, _ := queryBirdTableByName("crow")

	birds = append(birds, crow, jackdaw)
	return birds
}

func welcomePageHandler(w http.ResponseWriter, r *http.Request) {
	pagePath := "pages_html/welcome.html"

	tmpl, err := template.ParseFiles(pagePath)
	if err != nil {
		log.Printf("Error parsing html file in path %s\nerr: %s", pagePath, err)
	}

	page := WelcomePage{"Birdwiki", getBirdsArray()}

	tmpl.Execute(w, page)
}
