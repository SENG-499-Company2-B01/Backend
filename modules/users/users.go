package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// User represents a user entity
type User struct {
	Username       string            `json:"username"`
	Email          string            `json:"email"`
	Password       string            `json:"password"`
	Firstname      string            `json:"firstname"`
	LastName       string            `json:"lastname"`
	IsAdmin        bool              `json:"isAdmin"`
	Preferences    map[string]string `json:"preferences"`
	Qualifications []string          `json:"qualifications"`
}

// CreateUser handles the creation of a new user
func CreateUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("CreateUser function called.")

	// Parse request body into User struct
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		// If there is an error decoding the request body,
		// return a bad request response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		// If there is an error inserting the user into the collection,
		// log the error and return an internal server error response
		fmt.Println("Error inserting user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful user creation
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User created successfully")
}

// GetUsers retrieves all users from the MongoDB collection
func GetUsers(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("GetUsers function called.")

	// Define an empty slice to store the users
	var users []User

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// If there is an error retrieving users,
		// log the error and return an internal server error response
		fmt.Println("Error retrieving users:", err)
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
			fmt.Println("Error decoding user:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		// If there is an error iterating through the cursor,
		// log the error and return an internal server error response
		fmt.Println("Error iterating cursor:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response with the retrieved users
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUser retrieves a user by username
func GetUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("GetUser function called.")

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
			fmt.Println("User not found:", err)
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// If there is an error retrieving the user,
			// log the error and return an internal server error response
			fmt.Println("Error getting user:", err)
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
	fmt.Println("UpdateUser function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	username := strings.TrimPrefix(path, "/users/")
	username = strings.TrimSpace(username)

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

	// Update the user in the MongoDB collection
	filter := bson.M{"username": username}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// If there is an error updating the user in the collection,
		// log the error and return an internal server error response
		fmt.Println("Error updating user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful user update
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User updated successfully")
}

// DeleteUser handles the deletion of a user
func DeleteUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	fmt.Println("DeleteUser function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	username := strings.TrimPrefix(path, "/users/")
	username = strings.TrimSpace(username)

	// Delete the user from the MongoDB collection
	filter := bson.M{"username": username}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		// If there is an error deleting the user from the collection,
		// log the error and return an internal server error response
		fmt.Println("Error deleting user:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful user deletion
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User deleted successfully")
}
