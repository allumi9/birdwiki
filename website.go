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
	PicturePath string
}

type WelcomePage struct {
	Title string
	Birds []Bird
}

type SearchResultsPage struct {
	Title string
	Query string
	Birds []Bird
}

type BirdInfoPage struct {
	Title string
	Bird  Bird
}

func getBird(birdPath string, picPath string) (Bird, error) {
	fileDescription, err := os.ReadFile(birdPath)
	if err != nil {
		errorMessage := "Error when reading file" + birdPath
		log.Printf(errorMessage, err)
		return Bird{"", "", ""}, errors.New(errorMessage)
	}

	slicedDesc := strings.Split(string(fileDescription), "\n")
	return Bird{slicedDesc[0], slicedDesc[1], picPath}, nil
}

func getBirdsArray() []Bird {
	var birds []Bird

	jackdaw, _ := queryBirdTableByName("jackdaw")
	crow, _ := queryBirdTableByName("crow")

	birds = append(birds, crow, jackdaw)
	return birds
}

func getPreCheckedTemplate(pagePath string) *template.Template {
	tmpl, err := template.ParseFiles(pagePath)
	if err != nil {
		// I don't think having non existent pages is a good thing in the first place -> fatal
		closePostgresConnection(postgresConnection)
		log.Fatalf("Error parsing html file in path %s\nerr: %s", pagePath, err)
	}
	return tmpl
}

func welcomePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := getPreCheckedTemplate("pages_html/welcome.html")
	page := WelcomePage{"Birdwiki", getBirdsArray()}

	tmpl.Execute(w, page)
}
func returnHomePage(tmpl *template.Template, w http.ResponseWriter) {
	tmpl.Execute(w, WelcomePage{"Birdwiki", getBirdsArray()})
}

// TODO: replace home page returns with some like a bird not found page
func birdInfoPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := getPreCheckedTemplate("pages_html/bird_info.html")
	nameParam := r.URL.Query().Get("name")
	if len(clearString(strings.Trim(nameParam, " \n\r"))) == 0 {
		log.Println("Error: name parameter given to bird info page is empty or invalid")
		returnHomePage(tmpl, w)
		return
	}

	bird, err := queryBirdTableByName(nameParam)
	if err != nil {
		log.Println("Error when getting bird from database")
		returnHomePage(tmpl, w)
		return
	}
	log.Println(bird)
	page := BirdInfoPage{"Birdwiki", bird}

	tmpl.Execute(w, page)
}

func searchRequestHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form values: %v\n", err)
	}

	SearchQuery := r.FormValue("search-query")

	tmpl := getPreCheckedTemplate("pages_html/search_results.html")
	Birds, err := queryBirdArrayByName(SearchQuery)
	if err != nil {
		log.Println(err)
	}
	title := "q: " + SearchQuery
	page := SearchResultsPage{title, SearchQuery, Birds}

	tmpl.Execute(w, page)
}
