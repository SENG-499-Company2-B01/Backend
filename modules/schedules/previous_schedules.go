package schedules

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SENG-499-Company2-B01/Backend/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreatePrevSchedule handles the creation of a new past-schedule
func CreatePrevSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("CreatePrevSchedule function called.")

	// Parse request body into Schedule struct
	var newSchedule Schedule
	err := json.NewDecoder(r.Body).Decode(&newSchedule)
	if err != nil {
		// If there is an error decoding the request body,
		// log the error and return a bad request response
		logger.Error(fmt.Errorf("Error decoding the request body: "+err.Error()), http.StatusBadRequest)
		http.Error(w, "Error decoding the request body.", http.StatusBadRequest)
		return
	}

	// Insert the schedule into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newSchedule)
	if err != nil {
		// If there is an error inserting the schedule into the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error inserting schedule: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error inserting schedule.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful schedule creation
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Schedule created successfully")

	// Uncomment the follow line for debugging
	// logger.Info("CreateSchedule function completed.")
}

// GetSchedules retrieves all previous schedules from the MongoDB collection
func GetPrevSchedules(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetPrevSchedules function called.")

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
