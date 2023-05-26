package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// User represents a user entity
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// createUser handles the creation of a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body into User struct
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUser)
}

// getUser retrieves a user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the URL path
	userID := r.URL.Query().Get("id")

	// Retrieve the user from the MongoDB collection
	filter := bson.M{"id": userID}
	var user User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Send a response
	// ...
}

// updateUser handles updating an existing user
func updateUser(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the URL path
	userID := r.URL.Query().Get("id")

	// Parse request body into User struct
	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the user in the MongoDB collection
	filter := bson.M{"id": userID}
	update := bson.M{"$set": updatedUser}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

// deleteUser handles the deletion of a user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the URL path
	userID := r.URL.Query().Get("id")

	// Delete the user from the MongoDB collection
	filter := bson.M{"id": userID}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User deleted successfully")
}
