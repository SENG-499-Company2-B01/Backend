package users

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

// User represents a user entity
type User struct {
	Username       string            `json:"username" bson:"username"`
	Email          string            `json:"email" bson:"email"`
	Password       string            `json:"password" bson:"password"`
	Firstname      string            `json:"firstname" bson:"firstname"`
	LastName       string            `json:"lastname" bson:"lastname"`
	IsAdmin        bool              `json:"-" bson:"isAdmin"`
	Preferences    map[string]string `json:"preferences" bson:"preferences"`
	Qualifications []string          `json:"qualifications" bson:"qualifications"`
}

// CreateUser handles the creation of a new user
func CreateUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("CreateUser function called.")

	// Parse request body into User struct
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		// If there is an error decoding the request body,
		// log the error and return a bad request response
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// by default IsAdmin is supposed to be set to false
	newUser.IsAdmin = false

	// Insert the user into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		// If there is an error inserting the user into the collection,
		// log the error and return an internal server error response
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful user creation
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User created successfully")
}

// GetUsers retrieves all users from the MongoDB collection
func GetUsers(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetUsers function called.")

	// Define an empty slice to store the users
	var users []User

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// If there is an error retrieving users,
		// log the error and return an internal server error response
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into a User struct
	for cursor.Next(context.TODO()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			// If there is an error decoding a user document,
			// log the error and return an internal server error response
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		// If there is an error iterating through the cursor,
		// log the error and return an internal server error response
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response with the retrieved users
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUser retrieves a user by username
func GetUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetUser function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	username := strings.TrimPrefix(path, "/users/")
	username = strings.TrimSpace(username)

	// Retrieve the user from the MongoDB collection
	filter := bson.M{"username": username}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If the user is not found,
			// log the error and return a not found response
			logger.Error(fmt.Errorf("User not found"))
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// If there is an error retrieving the user,
			// log the error and return an internal server error response
			logger.Error(fmt.Errorf("Error getting user"))
			http.Error(w, "Error getting user", http.StatusInternalServerError)
		}
		return
	}

	// Send a response with the retrieved user
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles updating an existing user
func UpdateUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("UpdateUser function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	username := strings.TrimPrefix(path, "/users/")
	username = strings.TrimSpace(username)

	// Check if the user exists
	if !userExists(username, collection) {
		logger.Error(fmt.Errorf("User not found"))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse request body into a map
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		// If there is an error decoding the request body, log the error and return a bad request response
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// isAdmin cannot be updated
	if requestBody["isAdmin"] != nil {
		http.Error(w, "isAdmin Field cannot be updated", http.StatusInternalServerError)
		return
	}

	// Construct the update query
	update := bson.M{"$set": requestBody}

	// Update the user in the MongoDB collection
	filter := bson.M{"username": username}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// If there is an error updating the user in the collection,
		// log the error and return an internal server error response
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful user update
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User updated successfully")
}

// DeleteUser handles the deletion of a user
func DeleteUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("DeleteUser function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	username := strings.TrimPrefix(path, "/users/")
	username = strings.TrimSpace(username)

	// Check if the user exists
	if !userExists(username, collection) {
		logger.Error(fmt.Errorf("User not found"))
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete the user from the MongoDB collection
	filter := bson.M{"username": username}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		// If there is an error deleting the user from the collection,
		// log the error and return an internal server error response
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful user deletion
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User deleted successfully")
}

// Check if a user exists in the MongoDB collection
func userExists(username string, collection *mongo.Collection) bool {
	filter := bson.M{"username": username}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	return err == nil
}
