//code i used for demo
//P.S dont run this file

/*
type Article struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string             `json: "title,omitempty" bson:"title,omitempty"`
	Subtitle string             `json: "subtitle,omitempty" bson:"subtitle,omitempty"`
	Content  string             `json: "content,omitempty" bson:"content,omitempty"`

		ID       primitive.ObjectID `json:"Id"`
	Title    string `json: "title"`
	Subtitle string `json: "subtitle"`
	Content  string `json: "content"`
}
*/

/*
func returnArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("endpoint hit: returnArticle")

	json.NewEncoder(w).Encode(Articles)
}
*/

/*
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
*/
