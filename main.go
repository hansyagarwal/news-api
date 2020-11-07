package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the news api")
	fmt.Println("endpoint hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", home)

	http.HandleFunc("/articles", returnArticle)
	http.HandleFunc("/articles/1", returnById)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

type Article struct {
	Id       int    `json:"Id"`
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
	i := r.RequestURI
	a := r.GetBody
	if strings.HasPrefix(i, "/articles/") {
		pe := url.PathEscape(strings.TrimLeft(i, "/articles/"))
		fmt.Println(pe)
		fmt.Println("endpoint hit: returnById")
		json.NewEncoder(w).Encode(Articles)
	}
	fmt.Println("endpoint hit: returnArticle")

}

func main() {

	Articles = []Article{
		Article{Id: 1, Title: "Hello", Subtitle: "Sub hello", Content: "content1"},
		Article{Id: 2, Title: "Hello 2", Subtitle: "Sub hello 2", Content: "content2"},
	}

	handleRequests()
}
