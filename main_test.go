package main

import (
	"backend/tests/classrooms"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	// "encoding/json"
	// "fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/SENG-499-Company2-B01/Backend/tests/classrooms"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var router = mux.NewRouter()

func TestMain(m *testing.M) {
	setupRoutes(router)
	code := m.Run()
	os.Exit(code)
}

func TestAll(t *testing.T) {
	// Add all tests here
	classrooms.TestInsertClassroom(t, executeRequest, client)
	classrooms.TestGetClassroom(t, executeRequest, client)
	classrooms.TestGetClassrooms(t, executeRequest, client)
	classrooms.TestUpdateClassroom(t, executeRequest, client)
	classrooms.TestDeleteClassroom(t, executeRequest, client)

}

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

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
