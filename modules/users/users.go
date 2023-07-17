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
	Username      string                `json:"username" bson:"username"`
	Email         string                `json:"email" bson:"email"`
	Password      string                `json:"password" bson:"password"`
	Name          string                `json:"name" bson:"name"`
	IsAdmin       bool                  `json:"-" bson:"isAdmin"`
	Peng          bool                  `json:"peng" bson:"peng"`
	Pref_approved bool                  `json:"pref_approved" bson:"pref_approved"`
	Max_courses   int                   `json:"max_courses" bson:"max_courses"`
	Course_pref   []string              `json:"course_pref" bson:"course_pref"`
	Available     map[string][][]string `json:"available" bson:"available"`
}

func CreateUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("CreateUser function called.")

	// Parse request body into User struct
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		// If there is an error decoding the request body,
		// log the error and return a bad request response
		logger.Error(fmt.Errorf("Error decoding the request body: "+err.Error()), http.StatusBadRequest)
		http.Error(w, "Error decoding the request body.", http.StatusBadRequest)
		return
	}

	// Check if username or email already exists in the collection
	filter := bson.M{
		"$or": []bson.M{
			{"username": newUser.Username},
			{"email": newUser.Email},
		},
	}
	count, err := collection.CountDocuments(context.TODO(), filter, nil)
	if err != nil {
		// If there is an error querying the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error checking the collection: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error checking the collection.", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		// If the count is greater than 0, indicating an existing user,
		// return a conflict response
		logger.Error(fmt.Errorf("username or email already exists"), http.StatusConflict)
		http.Error(w, "Username or email already exists.", http.StatusConflict)
		return
	}

	// by default IsAdmin is supposed to be set to false
	newUser.IsAdmin = false

	// Insert the user into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		// If there is an error inserting the user into the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error inserting user: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error inserting user.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful user creation
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "User created successfully")

	// Uncomment the follow line for debugging
	// logger.Info("CreateUser function completed.")
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
		logger.Error(fmt.Errorf("Error retrieving users: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving users.", http.StatusInternalServerError)
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
			logger.Error(fmt.Errorf("Error decoding user document: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error decoding user document.", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		// If there is an error iterating through the cursor,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
		return
	}

	// Send a response with the retrieved users
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

	// Uncomment the follow line for debugging
	// logger.Info("GetUsers function completed.")
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
			logger.Error(fmt.Errorf("user not found"), http.StatusNotFound)
			http.Error(w, "User not found.", http.StatusNotFound)
		} else {
			// If there is an error retrieving the user,
			// log the error and return an internal server error response
			logger.Error(fmt.Errorf("error getting user"), http.StatusInternalServerError)
			http.Error(w, "Error getting user.", http.StatusInternalServerError)
		}
		return
	}

	// Send a response with the retrieved user
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

	// Uncomment the follow line for debugging
	// logger.Info("GetUser function completed.")
}

// UpdateUser handles updating an existing user
func UpdateUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("UpdateUser function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	username := strings.TrimPrefix(path, "/users/")
	username = strings.TrimSpace(username)

	// Check if the user exists in the collection
	filter := bson.M{"username": username}
	exists, err := userExists(filter, collection)
	if err != nil {
		// If there is an error querying the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error querying collection: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error querying collection.", http.StatusInternalServerError)
		return
	}
	if !exists {
		// If the semester doesn't exist,
		// return a not found response
		logger.Error(fmt.Errorf("user not found"), http.StatusNotFound)
		http.Error(w, "User not found.", http.StatusInternalServerError)
		return
	}

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

	// isAdmin cannot be updated
	if requestBody["isAdmin"] != nil {
		logger.Error(fmt.Errorf("isAdmin field cannot be updated"), http.StatusInternalServerError)
		http.Error(w, "isAdmin field cannot be updated.", http.StatusInternalServerError)
		return
	}

	// Construct the update query
	update := bson.M{"$set": requestBody}

	// Update the user in the MongoDB collection
	filter = bson.M{"username": username}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// If there is an error updating the user in the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error updating the user: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error updating the user.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful user update
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "User updated successfully.")

	// Uncomment the follow line for debugging
	// logger.Info("UpdateUser function completed.")
}

// userExists checks if a document exists in the collection based on a filter
func userExists(filter bson.M, collection *mongo.Collection) (bool, error) {
	count, err := collection.CountDocuments(context.TODO(), filter, nil)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
