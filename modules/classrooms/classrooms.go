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
	Shorthand   string   `json:"shorthand"`
	Building    string   `json:"building"`
	Capacity    int      `json:"capacity"`
	Room_number string   `json:"room_number"`
	Equipment   []string `json:"equipment"`
}

// CreateClassroom handles the creation of a new classroom
func CreateClassroom(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {

	// Parse r body into Classroom struct
	var newClassroom Classroom

	err := json.NewDecoder(r.Body).Decode(&newClassroom)
	if err != nil {
		// If there is an error decoding the r body,
		// return a bad r response
		http.Error(w, "Error decoding json", http.StatusBadRequest)
		return
	}

	filter := bson.M{"shorthand": newClassroom.Shorthand, "room_number": newClassroom.Room_number}
	// Store the retrieved classroom
	var existingClassroom Classroom
	// Find the classroom
	error := collection.FindOne(context.TODO(), filter).Decode(&existingClassroom)

	if error != nil && existingClassroom.Shorthand != "" {
		http.Error(w, "Error: Classroom already exists", http.StatusBadRequest)
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

func GetClassroom(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {

	// Parse r params
	vars := mux.Vars(r)
	shorthand, ok := vars["shorthand"]
	if !ok {
		fmt.Println("shorthand is missing in parameters")
		return
	}
	room_number, ok := vars["room_number"]
	if !ok {
		fmt.Println("room_number is missing in parameters")
		return
	}

	// Store the filter
	filter := bson.M{"shorthand": shorthand, "room_number": room_number}

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

func GetClassrooms(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {

	// Define an empty slice to store the classrooms
	var classrooms []Classroom

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// If there is an error retrieving classrooms,
		// log the error and return an internal server error response
		fmt.Println("Error retrieving classrooms:", err)
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

func UpdateClassroom(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {

	// Parse request params
	vars := mux.Vars(r)
	shorthand, ok := vars["shorthand"]
	if !ok {
		fmt.Println("shorthand is missing in parameters")
		return
	}
	room_number, ok := vars["room_number"]
	if !ok {
		fmt.Println("room_number is missing in parameters")
		return
	}

	// Store the filter
	filter := bson.M{"shorthand": shorthand, "room_number": room_number}

	// Parse r body into a map
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		// If there is an error decoding the r body,
		// return a bad r response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Construct the update query
	update := bson.M{"$set": requestBody}

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

// DeleteClassroom handles the deletion of a classroom
func DeleteClassroom(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {

	// Parse request params
	vars := mux.Vars(r)
	shorthand, ok := vars["shorthand"]
	if !ok {
		fmt.Println("shorthand is missing in parameters")
		return
	}
	room_number, ok := vars["room_number"]
	if !ok {
		fmt.Println("room_number is missing in parameters")
		return
	}

	// Store the filter
	filter := bson.M{"shorthand": shorthand, "room_number": room_number}

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
