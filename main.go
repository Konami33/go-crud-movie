package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Movie struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name   string             `json:"name"`
	Title  string             `json:"title"`
	Year   int                `json:"year"`
	Origin string             `json:"origin"`
}

var movies []Movie
var client *mongo.Client

func initMongoClient() {
	
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://konami9889:YDLk9ynz92ODfzsY@cluster0.2rqnsvs.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0YDLk9ynz92ODfzsY")

	opts = opts.SetServerAPIOptions(serverAPI)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func getMovies(w http.ResponseWriter, r *http.Request) {

	collection := client.Database("mydatabase").Collection("movies")
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var movie Movie
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("mydatabase").Collection("movies")
	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "Movie with ID %s has been deleted", id.Hex())
}

func getMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var movie Movie
	collection := client.Database("mydatabase").Collection("movies")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movie)
}

func createMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	collection := client.Database("mydatabase").Collection("movies")
	_, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var updatedMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&updatedMovie)

	collection := client.Database("mydatabase").Collection("movies")
	_, err := collection.ReplaceOne(context.Background(), bson.M{"_id": id}, updatedMovie)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(updatedMovie)

}

func main() {

	//db connection
	initMongoClient()
	defer client.Disconnect(context.Background())

	//router
	r := mux.NewRouter()

	//routes
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port http://localhost:8000/movies\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
