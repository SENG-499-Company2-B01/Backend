package main

import (
	"context"

	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SENG-499-Company2-B01/Backend/logger"
	"github.com/SENG-499-Company2-B01/Backend/modules/classrooms"
	"github.com/SENG-499-Company2-B01/Backend/modules/courses"
	"github.com/SENG-499-Company2-B01/Backend/modules/health"
	"github.com/SENG-499-Company2-B01/Backend/modules/middleware"
	"github.com/SENG-499-Company2-B01/Backend/modules/schedules"
	"github.com/SENG-499-Company2-B01/Backend/modules/users"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	// Create a "logs" directory
	logsDir := filepath.Join(dir, "logs")
	err = os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating logs directory:", err)
	}

	// Create a "logs.txt" file
	logPath := filepath.Join(logsDir, "logs.txt")

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}

	// Initialize the logger
	logger.InitLogger(os.Stdout, os.Stdout, os.Stderr, file)

	// Print a success message to the console
	logger.Info("Logger initialized successfully!")

	// Construct the path to the .env file
	envPath := filepath.Join(dir, ".env")

	// Load the .env file
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	if os.Getenv("ENVIRONMENT") == "development" {
		// Load the environment variables locally
		mongohost := os.Getenv("MONGO_LOCAL_HOST")
		mongoUsername := os.Getenv("MONGO_LOCAL_USERNAME")
		mongoPassword := os.Getenv("MONGO_LOCAL_PASSWORD")

		// Set up the MongoDB client with SCRAM-SHA-1 authentication
		clientOptions := options.Client().ApplyURI("mongodb://" + mongohost).
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

	} else if os.Getenv("ENVIRONMENT") == "production" {
		// Load the environment variables locally
		mongoUsername := os.Getenv("MONGO_PRODUCTION_USERNAME")
		mongoPassword := os.Getenv("MONGO_PRODUCTION_PASSWORD")
		mongoHost := os.Getenv("MONGO_PRODUCTION_HOST")

		// Use the MongoDB Atlas connection string
		connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", mongoUsername, mongoPassword, mongoHost)

		// Set up client options
		clientOptions := options.Client().ApplyURI(connectionString)

		// Connect to MongoDB
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Ping the MongoDB server to check the connection
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	logger.Info("Connected to MongoDB successfully!")
}

func handleUserRequests(router *mux.Router) {
	router.Use(func(next http.Handler) http.Handler {
		return middleware.Users_API_Access_Control(next, client.Database("schedule_db").Collection("users"))
	})

	// AUTHENTICATION
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		users.SignIn(w, r, client.Database("schedule_db").Collection("users"))
	}).Methods(http.MethodPost)

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		users.Logout(w, r, client.Database("schedule_db").Collection("users"))
	}).Methods(http.MethodPost)

	// Users CRUD Operations
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users.CreateUser(w, r, client.Database("schedule_db").Collection("users"))
	}).Methods(http.MethodPost)

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

}

func handleClassroomRequests(router *mux.Router) {
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
}

func handleCourseRequests(router *mux.Router) {
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
}

func handleScheduleRequests(router *mux.Router) {

	// Previous Schedule Operations
	router.HandleFunc("/schedules/prev", func(w http.ResponseWriter, r *http.Request) {
		schedules.GetSchedules(w, r, client.Database("schedule_db").Collection("previous_schedules"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/schedules/prev", func(w http.ResponseWriter, r *http.Request) {
		schedules.ApproveSchedule(w, r, client.Database("schedule_db").Collection("draft_schedules"), client.Database("schedule_db").Collection("previous_schedules"))
	}).Methods(http.MethodPost)

	// Schedules Read Operation
	router.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		schedules.GetSchedules(w, r, client.Database("schedule_db").Collection("draft_schedules"))
	}).Methods(http.MethodGet)

	// Schedules Generation Endpoints
	router.HandleFunc("/schedules/{year}/{term}/generate", func(w http.ResponseWriter, r *http.Request) {
		schedules.GenerateSchedule(w, r, client.Database("schedule_db").Collection("draft_schedules"))
	}).Methods(http.MethodPost)
}

func main() {
	// // Example Logging messages
	// logger.Info("This is an info message")
	// logger.Warning("This is a warning message")
	// logger.Error(fmt.Errorf("This is an error message"))

	router := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	handleUserRequests(router)
	handleClassroomRequests(router)
	handleCourseRequests(router)
	handleScheduleRequests(router)

	// This route will be used by the cloud server to test its health, it only ever returns 200 OK
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		health.CheckHealth(w, r)
	}).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

// // Example handle request
// router.HandleFunc("/", homePage).Methods(http.MethodGet)

// // Example Endpoint
// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the HomePage!")
// 	fmt.Println("Endpoint Hit: homePage")
// }
