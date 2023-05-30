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

// CreateSchedule handles the creation of a new schedule
func CreateClassroom(writer http.ResponseWriter, request *http.Request, collection *mongo.Collection) {
	fmt.Println("CreateClassroom function called.")

	// Parse request body into Classroom struct
	var newClassroom Classroom

	err := json.NewDecoder(request.Body).Decode(&newClassroom)
	fmt.Print(request.Body)
	if err != nil {
		// If there is an error decoding the request body,
		// return a bad request response
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the schedule into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newClassroom)
	if err != nil {
		// If there is an error inserting the classroom into the collection,
		// log the error and return an internal server error response
		fmt.Println("Error inserting schedule:", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful classroom creation
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Classroom created successfully")
}

func GetClassroom(writer http.ResponseWriter, request *http.Request, collection *mongo.Collection) {
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
			http.Error(writer, "Classroom not found", http.StatusNotFound)
		} else {
			// If there is an error retrieving the classroom,
			// log the error and return an internal server error response
			fmt.Println("Error getting classroom:", err)
			http.Error(writer, "Error getting classroom", http.StatusInternalServerError)
		}
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(classroom)
}

func GetClassroom(writer http.ResponseWriter, request *http.Request, collection *mongo.Collection) {
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
			http.Error(writer, "Classroom not found", http.StatusNotFound)
		} else {
			// If there is an error retrieving the classroom,
			// log the error and return an internal server error response
			fmt.Println("Error getting classroom:", err)
			http.Error(writer, "Error getting classroom", http.StatusInternalServerError)
		}
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(classroom)
}
