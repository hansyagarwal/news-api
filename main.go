package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func handleRequests() {
	articleHandlers := newArticleHandlers()

	http.HandleFunc("/", home)
	http.HandleFunc("/articles", articleHandlers.articles)
	http.HandleFunc("/articles/", ReturnById)
	http.HandleFunc("/articles/search", SearchArticle)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

type Article struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json: "title,omitempty" bson:"title,omitempty"`
	Subtitle  string             `json: "subtitle,omitempty" bson:"subtitle,omitempty"`
	Content   string             `json: "content,omitempty" bson:"content,omitempty"`
	CreatedAt time.Time          `json: "created_at,omitempty" bson:"created_at,omitempty"`
}

type articleHandlers struct {
	sync.Mutex
	store map[string]Article
}

var Articles []Article

//home function
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the news api")
	fmt.Println("endpoint hit: homePage")
}

//function to check if "/articles" request is GET or POST
func (h *articleHandlers) articles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

//function to handle GET request for listing all articles
func (h *articleHandlers) get(w http.ResponseWriter, r *http.Request) {
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

//function to handle POST request of creating article
func (h *articleHandlers) post(w http.ResponseWriter, r *http.Request) {
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

	result, err := articlesCollection.UpdateOne(ctx, bson.M{"title": a.Title}, bson.D{{"$set", bson.D{{"createdat", time.Now()}}}})
	_ = result
	fmt.Println(createResult.InsertedID)
}

func newArticleHandlers() *articleHandlers {
	return &articleHandlers{
		store: map[string]Article{},
	}
}

//function to handle GET request to get article by it's id
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

	result := Article{}
	collection := client.Database("news").Collection("articles")
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(w).Encode(result)
}

//function to handle GET request to search for an article by its title,subtitle or content
func SearchArticle(w http.ResponseWriter, r *http.Request) {
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

	keys, ok := r.URL.Query()["q"]

	if !ok || len(keys[0]) < 1 {
		log.Println("URL Param 'q' is missing")
		return
	}

	q := keys[0]
	log.Println("URL Parama 'key' is: " + string(q))
	q = string(q)
	result := Article{}
	collection := client.Database("news").Collection("articles")
	err = collection.FindOne(ctx, bson.M{"title": q}).Decode(&result)
	if err != nil {
		err = collection.FindOne(ctx, bson.M{"subtitle": q}).Decode(&result)
		if err != nil {
			err = collection.FindOne(ctx, bson.M{"content": q}).Decode(&result)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "` + err.Error() + `"}`))
				return
			}
		}
	}

	json.NewEncoder(w).Encode(result)
}

func main() {
	handleRequests()
}
