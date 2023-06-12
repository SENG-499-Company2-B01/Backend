package schedules

import (
	"context"
	"encoding/json"
	"strings"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Schedule represents a schedule entity
type Schedule struct {
	Semester 	string `json:"semester"`
	Classes     map[string]Class    `json:"classes"`
}

type Class struct {
	Course  string `json:"course"`
	Teacher string `json:"teacher"`
	Room    string `json:"room"`
}

// GenerateSchedule - Generates a new schedule
// TODO: Still needs to be implemented once algo team sets up their REST API
func GenerateSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("GenerateSchedule function called.")

	// Send a response indicating successful schedule creation
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "New Scheduled generated successfully")
}

// CreateSchedule handles the creation of a new schedule
func CreateSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("CreateSchedule function called.")

	// Parse request body into Schedule struct
	var newSchedule Schedule
	err := json.NewDecoder(r.Body).Decode(&newSchedule)
	if err != nil {
		// If there is an error decoding the request body,
		// return a bad request response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the schedule into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newSchedule)
	if err != nil {
		// If there is an error inserting the schedule into the collection,
		// log the error and return an internal server error response
		fmt.Println("Error inserting schedule:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful schedule creation
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Schedule created successfully")
}

// GetSchedules retrieves all schedules from the MongoDB collection
func GetSchedules(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("GetSchedules function called.")

	// Define an empty slice to store the schedules
	var schedules []Schedule

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// If there is an error retrieving schedules,
		// log the error and return an internal server error response
		fmt.Println("Error retrieving schedules:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			fmt.Println("Error decoding schedule:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		schedules = append(schedules, schedule)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		// If there is an error iterating through the cursor,
		// log the error and return an internal server error response
		fmt.Println("Error iterating cursor:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response with the retrieved schedules
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedules)
}

// GetSchedule retrieves a schedule by username
func GetSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("GetSchedule function called.")

	// Extract the schedule semester from the URL path
	path := r.URL.Path
	semester := strings.TrimPrefix(path, "/schedules/")
	semester = strings.TrimSpace(semester)

	// Retrieve the schedule from the MongoDB collection
	filter := bson.M{"semester": semester}
	var schedule Schedule
	err := collection.FindOne(context.TODO(), filter).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If the schedule is not found,
			// log the error and return a not found response
			fmt.Println("Schedule not found:", err)
			http.Error(w, "Schedule not found", http.StatusNotFound)
		} else {
			// If there is an error retrieving the schedule,
			// log the error and return an internal server error response
			fmt.Println("Error getting schedule:", err)
			http.Error(w, "Error getting schedule", http.StatusInternalServerError)
		}
		return
	}

	// Send a response with the retrieved schedule
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}

// UpdateSchedule handles updating an existing schedule
func UpdateSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("UpdateSchedule function called.")

	// Extract the schedule semester from the URL path
	path := r.URL.Path
	semester := strings.TrimPrefix(path, "/schedules/")
	semester = strings.TrimSpace(semester)

	// Parse request body into a map
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		// If there is an error decoding the request body,
		// return a bad request response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Construct the update query
	update := bson.M{"$set": requestBody}

	// Update the schedule in the MongoDB collection
	filter := bson.M{"semester": semester}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// If there is an error updating the schedule in the collection,
		// log the error and return an internal server error response
		fmt.Println("Error updating schedule:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful schedule update
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Schedule updated successfully")
}

// DeleteSchedule handles the deletion of a schedule
func DeleteSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("DeleteSchedule function called.")

	// Extract the schedule semester from the URL path
	path := r.URL.Path
	semester := strings.TrimPrefix(path, "/schedules/")
	semester = strings.TrimSpace(semester)

	// Delete the schedule from the MongoDB collection
	filter := bson.M{"semester": semester}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		// If there is an error deleting the schedule from the collection,
		// log the error and return an internal server error response
		fmt.Println("Error deleting schedule:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful schedule deletion
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Schedule deleted successfully")
}