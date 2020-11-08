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
	//fmt.Println("endpoint hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", home)

	//http.HandleFunc("/articles", returnArticle)
	http.HandleFunc("/articles/", returnById)
	//http.HandleFunc("/articles", CreateArticle)
	http.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			//CreateArticle()
		} else {
			//returnArticle()
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}

type Article struct {
	Id       string `json:"Id"`
	Title    string `json: "title"`
	Subtitle string `json: "subtitle"`
	Content  string `json: "content"`
	//CreationTime
}

/*
type Article struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string             `json: "title,omitempty" bson:"title,omitempty"`
	Subtitle string             `json: "subtitle,omitempty" bson:"subtitle,omitempty"`
	Content  string             `json: "content,omitempty" bson:"content,omitempty"`
}
*/

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
	Articles = []Article{
		Article{Id: "1", Title: "Hello", Subtitle: "Sub hello", Content: "content1"},
		Article{Id: "2", Title: "Hello 2", Subtitle: "Sub hello 2", Content: "content2"},
	}

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

	articleResult, err := articlesCollection.InsertMany(ctx, []interface{}{
		bson.D{
			{Key: "id", Value: "1"},
			{Key: "title", Value: "hello"},
			{Key: "subtitle", Value: "hello sub"},
			{Key: "content", Value: "content"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inseted %v documents into article collection!\n", len(articleResult.InsertedIDs))
	//err = client.Ping(ctx, readpref.Primary())
	//if err != nil {
	//	log.Fatal(err)
	//}

	handleRequests()
}
