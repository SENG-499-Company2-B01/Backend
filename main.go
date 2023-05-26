package main

import (
	"context"
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

func init() {
	// Set up the MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Set the database and collection
	database := client.Database("schedule_db")
	collection = database.Collection("users")
}

func handleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/users", createUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", getUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", updateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", deleteUser).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
