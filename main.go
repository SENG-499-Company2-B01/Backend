package main

import (
	"context"

	"log"
	// "fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SENG-499-Company2-B01/Backend/logger"
	"github.com/SENG-499-Company2-B01/Backend/modules/classrooms"
	"github.com/SENG-499-Company2-B01/Backend/modules/courses"
	"github.com/SENG-499-Company2-B01/Backend/modules/schedules"
	"github.com/SENG-499-Company2-B01/Backend/modules/users"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	// Initialize the logger
	var err error
	err = logger.InitLogger(os.Stdout, os.Stdout, os.Stderr)
	if err != nil {
		// log.Fatal("Error initializing logger:", err)
		logger.Error(err)
		return
	}

	// Print a success message to the console
	logger.Info("Logger initialized successfully!")

	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory:", err)
	}

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

	logger.Info("Connected to MongoDB successfully!")
}

func handleRequests() {
	router := mux.NewRouter()

	// // Example handle request
	// router.HandleFunc("/", homePage).Methods(http.MethodGet)

	// Courses CRUD Operations
	router.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		courses.CreateCourse(w, r, client.Database("schedule_db").Collection("courses"))
	}).Methods(http.MethodPost)

	router.HandleFunc("/courses/{courseShortHand}", func(w http.ResponseWriter, r *http.Request) {
		courses.GetCourse(w, r, client.Database("schedule_db").Collection("courses"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		courses.GetCourses(w, r, client.Database("schedule_db").Collection("courses"))
	}).Methods(http.MethodGet)
	router.HandleFunc("/courses/{courseShortHand}", func(w http.ResponseWriter, r *http.Request) {
		courses.DeleteCourse(w, r, client.Database("schedule_db").Collection("courses"))
	}).Methods(http.MethodDelete)

	router.HandleFunc("/courses/{courseShortHand}", func(w http.ResponseWriter, r *http.Request) {
		courses.UpdateCourse(w, r, client.Database("schedule_db").Collection("courses"))
	}).Methods(http.MethodPut)

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

	// Classroom CRUD Operations
	router.HandleFunc("/classrooms", func(w http.ResponseWriter, r *http.Request) {
		classrooms.CreateClassroom(w, r, client.Database("schedule_db").Collection("classrooms"))
	}).Methods(http.MethodPost)

	router.HandleFunc("/classrooms/{shorthand}/{room_number}", func(w http.ResponseWriter, r *http.Request) {
		classrooms.GetClassroom(w, r, client.Database("schedule_db").Collection("classrooms"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/classrooms", func(w http.ResponseWriter, r *http.Request) {
		classrooms.GetClassrooms(w, r, client.Database("schedule_db").Collection("classrooms"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/classrooms/{shorthand}/{room_number}", func(w http.ResponseWriter, r *http.Request) {
		classrooms.UpdateClassroom(w, r, client.Database("schedule_db").Collection("classrooms"))
	}).Methods(http.MethodPut)

	router.HandleFunc("/classrooms/{shorthand}/{room_number}", func(w http.ResponseWriter, r *http.Request) {
		classrooms.DeleteClassroom(w, r, client.Database("schedule_db").Collection("classrooms"))
	}).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8000", router))
}

// // Example Endpoint
// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the HomePage!")
// 	fmt.Println("Endpoint Hit: homePage")
// }

func main() {
	// // Example Logging messages
	// logger.Info("This is an info message")
	// logger.Warning("This is a warning message")
	// logger.Error(fmt.Errorf("This is an error message"))
	

	handleRequests()
}
