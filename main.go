package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the news api")
	//fmt.Println("endpoint hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", home)

	http.HandleFunc("/articles", returnArticle)
	http.HandleFunc("/articles/", returnById)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

type Article struct {
	Id       string `json:"Id"`
	Title    string `json: "title"`
	Subtitle string `json: "subtitle"`
	Content  string `json: "content"`
	//CreationTime
}

var Articles []Article

func returnArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: returnArticle")
	json.NewEncoder(w).Encode(Articles)
}

func returnById(w http.ResponseWriter, r *http.Request) {
	//i := r.RequestURI

	id := strings.TrimPrefix(r.URL.Path, "/articles/")
	fmt.Println(id)
	for _, article := range Articles {
		if article.Id == id {
			json.NewEncoder(w).Encode(article)
		}
	}
	/*
		fmt.Println(i)
		if strings.HasPrefix(i, "/articles/") {
			key := url.PathEscape(strings.TrimLeft(i, "/articles/"))

			fmt.Println(key)
			fmt.Println("endpoint hit: returnById")
			for _, article := range Articles {
				if article.Id == key {
					json.NewEncoder(w).Encode(article)
				}
			}
		}
	*/
}

func main() {

	Articles = []Article{
		Article{Id: "1", Title: "Hello", Subtitle: "Sub hello", Content: "content1"},
		Article{Id: "2", Title: "Hello 2", Subtitle: "Sub hello 2", Content: "content2"},
	}

	handleRequests()
}
