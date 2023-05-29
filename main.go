package main

import (

	// "encoding/json"
	// "fmt"
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/bson"

	c "main_api/Modules/Courses"
)

var client *mongo.Client

func initDB() {
	// Get the current working directory
	var err error
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory:", err)
	}
	log.Println("Current working directory:", dir)

	// Construct the path to the .env file
	envPath := filepath.Join(dir, ".env")

	// Load the .env file
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Load the environment variables locally
	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoAddress := os.Getenv("MONGO_ADDRESS")
	mongoPort := os.Getenv("MONGO_PORT")

	// Set up the MongoDB client with SCRAM-SHA-1 authentication
	clientOptions := options.Client().ApplyURI("mongodb://" + mongoAddress + ":" + mongoPort).
		SetAuth(options.Credential{
			Username:      mongoUsername,
			Password:      mongoPassword,
			AuthMechanism: "SCRAM-SHA-256",
		})

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
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)

	cm := c.CourseModule{}
	router.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		cm.CreateCourse(w, r, client.Database("schedule_db").Collection("courses"))
	}).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	initDB()
	handleRequests()
}
