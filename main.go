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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the news api")
	fmt.Println("endpoint hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", home)

	http.HandleFunc("/articles", ReturnAllArticles)
	http.HandleFunc("/articles/", ReturnById)
	http.HandleFunc("/articless", CreateArticle)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type Article struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string             `json: "title,omitempty" bson:"title,omitempty"`
	Subtitle string             `json: "subtitle,omitempty" bson:"subtitle,omitempty"`
	Content  string             `json: "content,omitempty" bson:"content,omitempty"`
	//CreationTime
}

var Articles []Article

func ReturnAllArticles(w http.ResponseWriter, r *http.Request) {
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

func ReturnById(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("content-type", "application/json")

	params := strings.TrimPrefix(r.URL.Path, "/articles/")
	fmt.Println(params)
	docID, err := primitive.ObjectIDFromHex(params)

	//id, _ := primitive.ObjectIDFromHex(params["id"])
	//var article Article
	result := Article{}
	collection := client.Database("news").Collection("articles")
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(result)

	/*
		for _, article := range Articles {
			if article.Id == id {
				json.NewEncoder(w).Encode(article)
			}
		}*/
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
