package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the news api")
	fmt.Println("endpoint hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", home)

	http.HandleFunc("/articles", returnAllArticles)
	http.HandleFunc("/articles/", returnById)
	http.HandleFunc("/articless", CreateArticle)
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

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://news-api:qwertyuiop@cluster0.gv8ol.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	collection := client.Database("news").Collection("articles")
	var arti []Article
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var art Article
		cursor.Decode(&art)
		arti = append(arti, art)
	}

	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(arti)
}

func returnById(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/articles/")
	fmt.Println(id)
	for _, article := range Articles {
		if article.Id == id {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a Article
	err := decoder.Decode(&a)
	if err != nil {
		panic(err)
	}
	log.Println(a)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://news-api:qwertyuiop@cluster0.gv8ol.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	quickstartDb := client.Database("news")
	articlesCollection := quickstartDb.Collection("articles")

	createResult, err := articlesCollection.InsertOne(ctx, a)
	if err != nil {
		panic(err)
	}
	fmt.Println(createResult.InsertedID)
}

func main() {
	handleRequests()
}
