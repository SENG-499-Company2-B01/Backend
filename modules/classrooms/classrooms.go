package classrooms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Classroom struct {
	Id                 string
	Building           string
	Capacity           int
	RoomNumber         string
	AvailableEquipment []string
}

// CreateSchedule handles the creation of a new classroom
func CreateClassroom(w http.ResponseWriter, request *http.Request, collection *mongo.Collection) {
	fmt.Println("CreateClassroom function called.")

	// Parse request body into Classroom struct
	var newClassroom Classroom

	err := json.NewDecoder(request.Body).Decode(&newClassroom)
	fmt.Print(request.Body)
	if err != nil {
		// If there is an error decoding the request body,
		// return a bad request response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the classroom into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newClassroom)
	if err != nil {
		// If there is an error inserting the classroom into the collection,
		// log the error and return an internal server error response
		fmt.Println("Error inserting classroom:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful classroom creation
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Classroom created successfully")
}

func GetClassroom(w http.ResponseWriter, request *http.Request, collection *mongo.Collection) {
	fmt.Println("GetClassroom function called.")

	// Parse request params
	vars := mux.Vars(request)
	id, ok := vars["classroom"]
	if !ok {
		fmt.Println("classroom is missing in parameters")
		return
	}
	// Store the filter
	filter := bson.M{"id": id}

	// Store the retrieved classroom
	var classroom Classroom

	// Find the classroom
	err := collection.FindOne(context.TODO(), filter).Decode(&classroom)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If the classroom is not found,
			// log the error and return a not found response
			fmt.Println("Classroom not found:", err)
			http.Error(w, "Classroom not found", http.StatusNotFound)
		} else {
			// If there is an error retrieving the classroom,
			// log the error and return an internal server error response
			fmt.Println("Error getting classroom:", err)
			http.Error(w, "Error getting classroom", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(classroom)
}

func GetClassrooms(w http.ResponseWriter, request *http.Request, collection *mongo.Collection) {
	fmt.Println("GetClassrooms function called.")

	// Define an empty slice to store the schedules
	var classrooms []Classroom

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// If there is an error retrieving classrooms,
		// log the error and return an internal server error response
		fmt.Println("Error retrieving schedules:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into a Classroom struct
	for cursor.Next(context.TODO()) {
		var classroom Classroom
		err := cursor.Decode(&classroom)
		if err != nil {
			// If there is an error decoding a classroom document,
			// log the error and return an internal server error response
			fmt.Println("Error decoding classroom:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		classrooms = append(classrooms, classroom)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		// If there is an error iterating through the cursor,
		// log the error and return an internal server error response
		fmt.Println("Error iterating cursor:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(classrooms)
}

func UpdateClassroom(w http.ResponseWriter, request *http.Request, collection *mongo.Collection) {
	fmt.Println("UpdateClassrooms function called.")

	// Parse request params
	vars := mux.Vars(request)
	id, ok := vars["classroom"]
	if !ok {
		fmt.Println("classroom is missing in parameters")
		return
	}

	// Parse request body into a map
	var requestBody map[string]interface{}
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if err != nil {
		// If there is an error decoding the request body,
		// return a bad request response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Construct the update query
	update := bson.M{"$set": requestBody}

	// Update the classroom in the MongoDB collection
	filter := bson.M{"id": id}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// If there is an error updating the classroom in the collection,
		// log the error and return an internal server error response
		fmt.Println("Error updating classroom:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful classroom update
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Classroom updated successfully")
}

// DeleteSchedule handles the deletion of a classroom
func DeleteClassroom(w http.ResponseWriter, request *http.Request, collection *mongo.Collection) {
	fmt.Println("DeleteClassroom function called.")

	// Parse request params
	vars := mux.Vars(request)
	id, ok := vars["classroom"]
	if !ok {
		fmt.Println("classroom is missing in parameters")
		return
	}

	// Delete the classroom from the MongoDB collection
	filter := bson.M{"id": id}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		// If there is an error deleting the classroom from the collection,
		// log the error and return an internal server error response
		fmt.Println("Error deleting classroom:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful classroom deletion
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Classroom deleted successfully")
}
