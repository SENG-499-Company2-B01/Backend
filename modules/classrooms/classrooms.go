package classrooms

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SENG-499-Company2-B01/Backend/logger"
)

type Classroom struct {
	Building    string `json:"building"`
	Capacity    int    `json:"capacity"`
	Room_number string `json:"room" bson:"room"`
}

// CreateClassroom handles the creation of a new classroom
func CreateClassroom(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("CreateClassroom function called.")

	// Parse r body into Classroom struct
	var newClassroom Classroom
	err := json.NewDecoder(r.Body).Decode(&newClassroom)
	if err != nil {
		// If there is an error decoding the request body,
		// log the error and return a bad request response
		logger.Error(fmt.Errorf("Error decoding the request body: "+err.Error()), http.StatusBadRequest)
		http.Error(w, "Error decoding the request body.", http.StatusBadRequest)
		return
	}

	// Check if shorthand already exists in the collection
	filter := bson.M{"building": newClassroom.Building}
	count, err := collection.CountDocuments(context.TODO(), filter, nil)

	if err != nil {
		// If there is an error querying the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error checking the collection: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error checking the collection.", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		// If the count is greater than 0, indicating an existing classroom,
		// return a conflict response
		logger.Error(fmt.Errorf("classroom already exists"), http.StatusConflict)
		http.Error(w, "Classroom already exists.", http.StatusConflict)
		return
	}

	// Insert the classroom into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newClassroom)
	if err != nil {
		// If there is an error inserting the classroom into the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error inserting classroom: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving classrooms.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful classroom creation
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Classroom created successfully")

	// Uncomment the follow line for debugging
	// logger.Info("CreateClassroom function completed.")
}

func GetClassrooms(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetClassrooms function called.")

	// Define an empty slice to store the classrooms
	var classrooms []Classroom

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// If there is an error retrieving classrooms,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error retrieving classrooms: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving classrooms.", http.StatusInternalServerError)
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
			logger.Error(fmt.Errorf("Error decoding classroom: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error decoding classroom.", http.StatusInternalServerError)
			return
		}
		classrooms = append(classrooms, classroom)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		// If there is an error iterating through the cursor,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(classrooms)

	// Uncomment the follow line for debugging
	// logger.Info("GetClassrooms function completed.")
}

func GetClassroom(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetClassroom function called.")

	// Parse request params
	vars := mux.Vars(r)
	building, ok := vars["building"]
	if !ok {
		logger.Error(fmt.Errorf("building is missing in parameters"), http.StatusBadRequest)
		return
	}
	room_number, ok := vars["room"]
	if !ok {
		logger.Error(fmt.Errorf("room is missing in parameters"), http.StatusBadRequest)
		return
	}

	// Store the filter
	filter := bson.M{"building": building, "room": room_number}

	// Store the retrieved classroom
	var classroom Classroom

	// Find the classroom
	err := collection.FindOne(context.TODO(), filter).Decode(&classroom)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If the classroom is not found,
			// log the error and return a not found response
			logger.Error(fmt.Errorf("Classroom not found: "+err.Error()), http.StatusNotFound)
			http.Error(w, "Classroom not found.", http.StatusNotFound)
		} else {
			// If there is an error retrieving the classroom,
			// log the error and return an internal server error response
			logger.Error(fmt.Errorf("Error getting classroom: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error getting classroom.", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(classroom)

	// Uncomment the follow line for debugging
	// logger.Info("GetClassroom function completed.")
}

func UpdateClassroom(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("UpdateClassroom function called.")

	// Parse request params
	vars := mux.Vars(r)
	building, ok := vars["building"]
	if !ok {
		logger.Error(fmt.Errorf("building is missing in parameters"), http.StatusBadRequest)
		return
	}
	room_number, ok := vars["room"]
	if !ok {
		logger.Error(fmt.Errorf("room is missing in parameters"), http.StatusBadRequest)
		return
	}

	// Check if the shorthand exists in the collection
	filter := bson.M{
		"$and": []bson.M{
			{"building": building},
			{"room": room_number},
		},
	}
	exists, err := classroomExists(filter, collection)
	if err != nil {
		// If there is an error querying the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error querying collection: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error querying collection.", http.StatusInternalServerError)
		return
	}
	if !exists {
		// If the classroom doesn't exist,
		// return a not found response
		logger.Error(fmt.Errorf("classroom not found"), http.StatusInternalServerError)
		http.Error(w, "Classroom not found.", http.StatusInternalServerError)
		return
	}

	// Store the filter
	filter = bson.M{"building": building, "room": room_number}

	// Parse request body into a map
	var requestBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		// If there is an error decoding the request body,
		// log the error and return a bad request response
		logger.Error(fmt.Errorf("Error decoding the request body: "+err.Error()), http.StatusBadRequest)
		http.Error(w, "Error decoding the request body.", http.StatusBadRequest)

		return
	}

	// Construct the update query
	update := bson.M{"$set": requestBody}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// If there is an error updating the classroom in the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error updating classroom: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error updating classroom.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful classroom update
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Classroom updated successfully")

	// Uncomment the follow line for debugging
	// logger.Info("UpdateClassroom function completed.")
}

// DeleteClassroom handles the deletion of a classroom
func DeleteClassroom(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("DeleteClassroom function called.")

	// Parse request params
	vars := mux.Vars(r)
	building, ok := vars["building"]
	if !ok {
		logger.Error(fmt.Errorf("building is missing in parameters"), http.StatusBadRequest)
		return
	}
	room_number, ok := vars["room"]
	if !ok {
		logger.Error(fmt.Errorf("room is missing in parameters"), http.StatusBadRequest)
		return
	}

	// Check if the shorthand exists in the collection
	filter := bson.M{
		"$and": []bson.M{
			{"building": building},
			{"room": room_number},
		},
	}
	exists, err := classroomExists(filter, collection)
	if err != nil {
		// If there is an error querying the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error querying collection: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error querying collection.", http.StatusInternalServerError)
		return
	}
	if !exists {
		// If the classroom doesn't exist,
		// return a not found response
		logger.Error(fmt.Errorf("classroom not found"), http.StatusInternalServerError)
		http.Error(w, "Classroom not found.", http.StatusInternalServerError)
		return
	}

	// Store the filter
	filter = bson.M{"building": building, "room": room_number}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		// If there is an error deleting the classroom from the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error deleting classroom: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error deleting classroom.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful classroom deletion
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Classroom deleted successfully")

	// Uncomment the follow line for debugging
	// logger.Info("DeleteClassroom function completed.")
}

// classroomExists checks if a document exists in the collection based on a filter
func classroomExists(filter bson.M, collection *mongo.Collection) (bool, error) {
	count, err := collection.CountDocuments(context.TODO(), filter, nil)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
