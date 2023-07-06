package schedules

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SENG-499-Company2-B01/Backend/logger"
)

// Schedule represents a schedule entity
type Schedule struct {
	Year  string `json:"year"`
	Terms []Term `json:"terms"`
}

type Term struct {
	Term    string           `json:"term"`
	Courses []CourseOffering `json:"courses"`
}

type CourseOffering struct {
	Course   string  `json:"course"`
	Sections []Class `json:"sections"`
}

type Class struct {
	Num       string   `json:"num"`
	Building  string   `json:"building"`
	Professor string   `json:"professor"`
	Days      []string `json:"days"`
	NumSeats  int      `json:"num_seats"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
}

// GenerateSchedule - Generates a new schedule
// TODO: Still needs to be implemented once algo team sets up their REST API
func GenerateSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GenerateSchedule function called.")

	// Extract the year and term values from the URL path
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 {
		// If there is an error parsing the url path,
		// log the error and return a bad request response
		logger.Error(fmt.Errorf("invalid URL path"), http.StatusBadRequest)
		http.Error(w, "Invalid URL path.", http.StatusBadRequest)
		return
	}

	// Extract the year and term from path
	// year := path[2]
	// term := path[3]

	// TODO: finish this function

	// Send a response indicating successful schedule creation
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("New Scheduled generated successfully")
	fmt.Fprintf(w, "New Scheduled generated successfully")

	// Uncomment the follow line for debugging
	// logger.Info("GenerateSchedule function completed.")
}

// ApproveSchedule - removes schedule in draft collection and adds it to previous_schedules collection, approving it.
func ApproveSchedule(w http.ResponseWriter, r *http.Request, draftsCollection *mongo.Collection, previousSchedulesCollection *mongo.Collection) {
	logger.Info("ApproveSchedule function called.")

	// Find the schedule in the "draft_schedule" collection
	var foundSchedule Schedule
	err := draftsCollection.FindOne(context.TODO(), bson.M{}).Decode(&foundSchedule)
	if err != nil {
		logger.Error(fmt.Errorf("failed to find schedule in drafts collection"), http.StatusInternalServerError)
		http.Error(w, "Failed to find schedule.", http.StatusInternalServerError)
		return
	}

	// Insert the found schedule into the "previous_schedules" collection
	_, err = previousSchedulesCollection.InsertOne(context.TODO(), foundSchedule)
	if err != nil {
		logger.Error(fmt.Errorf("failed to insert schedule into previous_schedules collection"), http.StatusInternalServerError)
		http.Error(w, "Failed to insert schedule.", http.StatusInternalServerError)
		return
	}

	// Delete the schedule from the "draft_schedule" collection
	_, err = draftsCollection.DeleteOne(context.TODO(), bson.M{})
	if err != nil {
		logger.Error(fmt.Errorf("failed to delete from drafts collection"), http.StatusInternalServerError)
		http.Error(w, "Failed to delete schedule.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful schedule creation
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Schedule has been approved")

	// Uncomment the follow line for debugging
	// logger.Info("ApproveSchedule function completed.")
}

// GetSchedules retrieves all schedules from the MongoDB collection
func GetSchedules(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetSchedules function called.")

	// Define an empty slice to store the schedules
	var schedules []Schedule

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// If there is an error retrieving schedules,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error retrieving schedules: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving schedules.", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into a Schedule struct
	for cursor.Next(context.TODO()) {
		var schedule Schedule
		err := cursor.Decode(&schedule)
		if err != nil {
			// If there is an error decoding a schedule document,
			// log the error and return an internal server error response
			logger.Error(fmt.Errorf("Error decoding schedule: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error decoding schedule.", http.StatusInternalServerError)
			return
		}
		schedules = append(schedules, schedule)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		// If there is an error iterating through the cursor,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
		return
	}

	// Send a response with the retrieved schedules
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedules)

	// Uncomment the follow line for debugging
	// logger.Info("GetSchedules function completed.")
}
