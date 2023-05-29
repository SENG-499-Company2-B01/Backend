package main

import (
	"context"
	"log"
	"os"
	"net/http"
	"path/filepath"

	"backend/modules/users"
	"backend/modules/schedules"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
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
			Username:   mongoUsername,
			Password:   mongoPassword,
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

	// // Example handle request
	// router.HandleFunc("/", homePage).Methods(http.MethodGet)

	// Users CRUD Operations
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users.CreateUser(w, r, client.Database("schedule_db").Collection("users"))
	}).Methods(http.MethodPost)

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users.GetUsers(w, r, client.Database("schedule_db").Collection("users"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/users/{username}", func(w http.ResponseWriter, r *http.Request) {
		users.GetUser(w, r, client.Database("schedule_db").Collection("users"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/users/{username}", func(w http.ResponseWriter, r *http.Request) {
		users.UpdateUser(w, r, client.Database("schedule_db").Collection("users"))
	}).Methods(http.MethodPut)

	router.HandleFunc("/users/{username}", func(w http.ResponseWriter, r *http.Request) {
		users.DeleteUser(w, r, client.Database("schedule_db").Collection("users"))
	}).Methods(http.MethodDelete)

	// Schedules CRUD Operations
	router.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		schedules.CreateSchedule(w, r, client.Database("schedule_db").Collection("schedules"))
	}).Methods(http.MethodPost)

	router.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		schedules.GetSchedules(w, r, client.Database("schedule_db").Collection("schedules"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/schedules/{schedule}", func(w http.ResponseWriter, r *http.Request) {
		schedules.GetSchedule(w, r, client.Database("schedule_db").Collection("schedules"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/schedules/{schedule}", func(w http.ResponseWriter, r *http.Request) {
		schedules.UpdateSchedule(w, r, client.Database("schedule_db").Collection("schedules"))
	}).Methods(http.MethodPut)

	router.HandleFunc("/schedules/{schedule}", func(w http.ResponseWriter, r *http.Request) {
		schedules.DeleteSchedule(w, r, client.Database("schedule_db").Collection("schedules"))
	}).Methods(http.MethodDelete)

	// Courses CRUD Operations


	// Classroom CRUD Operations


	log.Fatal(http.ListenAndServe(":8000", router))
}

// // Example Endpoint
// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the HomePage!")
// 	fmt.Println("Endpoint Hit: homePage")
// }

func main() {
	handleRequests()
}
