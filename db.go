package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Post struct {
	Title string `json:”title,omitempty”`

	Body string `json:”body,omitempty”`
}

func InsertPost(title string, body string) {

	post := Post{title, body}
	collection := client.Database("news").Collection("articles")
	insertResult, err := collection.InsertOne(context.TODO(), post)
	if err != nil {

		log.Fatal(err)

	}
	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}

func GetPost(id bson.ObjectId) {

	collection := client.Database("news").Collection("articles")

	filter := bson.D

	var post Post

	err := collection.FindOne(context.TODO(), filter).Decode(&post)

	if err != nil {

		log.Fatal(err)

	}

	fmt.Println("Found post with title ", post.Title)

}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://news-api:qwertyuiop@cluster0.gv8ol.mongodb.net/test"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

}
