package courses

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SENG-499-Company2-B01/Backend/logger"
)

type Course struct {
	ShortHand     string     `json:"shorthand" bson:"shorthand"`
	Name          string     `json:"name" bson:"name"`
	Equipements   []string   `json:"equipment" bson:"equipements"`
	Prerequisites [][]string `json:"prerequisites" bson:"prerequisites"`
}

// CreateCourse - creates a new course in the course DB
func CreateCourse(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("CreateCourse function called.")

	var newCourse Course
	err := json.NewDecoder(r.Body).Decode(&newCourse)
	if err != nil {
		// If there is an error decoding the request body,
		// log the error and return a bad request response
		logger.Error(fmt.Errorf("Error decoding the request body: "+err.Error()), http.StatusBadRequest)
		http.Error(w, "Error decoding the request body.", http.StatusBadRequest)
		return
	}

	// CHECK if shorthand is ABC101 format
	if !hasThreeConsecutiveNumerics(newCourse.ShortHand) {
		logger.Error(fmt.Errorf("invalid Course shorthand"), http.StatusBadRequest)
		http.Error(w, "Invalid Course shorthand.", http.StatusBadRequest)
		return
	}

	// CHECK if course doesn't exist
	var result Course
	filter := bson.D{{"shorthand", newCourse.ShortHand}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil && result.ShortHand != "" {
		logger.Error(fmt.Errorf("Course %s already exists", result.ShortHand), http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("Error: %s course already exists.", result.ShortHand), http.StatusInternalServerError)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		logger.Error(fmt.Errorf("Error while checking for duplicate document: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error while checking for duplicate document.", http.StatusInternalServerError)
		return
	}

	// Insert the user into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newCourse)
	if err != nil {
		logger.Error(fmt.Errorf("Error while inserting course into DB: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error while inserting course into DB.", http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newCourse)

	// Uncomment the follow line for debugging
	// logger.Info("CreateCourse function completed.")
}

// GetCourses - retieves all the courses from the DB
func GetCourses(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetCourses function called.")

	// Define an empty slice to store the users
	var courses []Course

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		logger.Error(fmt.Errorf("Error retrieving users: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving users.", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into a User struct
	for cursor.Next(context.TODO()) {
		var course Course
		err := cursor.Decode(&course)
		if err != nil {
			logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
			return
		}
		courses = append(courses, course)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
		return
	}

	// Send a response with the retrieved users
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)

	// Uncomment the follow line for debugging
	// logger.Info("GetCourses function completed.")
}

// GetCourse - gets course with the given course shorthand
func GetCourse(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetCourse function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	courseShortHand := strings.TrimPrefix(path, "/courses/")
	// CHECK if shorthand is ABC101 format
	if !hasThreeConsecutiveNumerics(courseShortHand) {
		logger.Error(fmt.Errorf("invalid Course shorthand"), http.StatusBadRequest)
		http.Error(w, "Invalid Course shorthand.", http.StatusBadRequest)
		return
	}

	// Get course from DB
	var result Course
	filter := bson.D{{"shorthand", courseShortHand}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		logger.Error(fmt.Errorf("Error while finding the course: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error while finding the course.", http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

	// Uncomment the follow line for debugging
	// logger.Info("GetCourse function completed.")
}

// UpdateCourse - update the course witht the given shorthand
func UpdateCourse(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("UpdateCourse function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	courseShortHand := strings.TrimPrefix(path, "/courses/")

	// Check if course exists
	var result Course
	filter := bson.D{{"shorthand", courseShortHand}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		logger.Error(fmt.Errorf("Course %s doesn't exist", courseShortHand), http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("Error: %s course doesn't exist.", courseShortHand), http.StatusInternalServerError)
		return
	}
	if err != nil {
		logger.Error(fmt.Errorf("Error while finding the course: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error while finding the course.", http.StatusInternalServerError)
		return
	}

	// extract body
	var updateCourse map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updateCourse)
	if err != nil {
		// If there is an error decoding the request body,
		// log the error and return a bad request response
		logger.Error(fmt.Errorf("Error decoding the request body: "+err.Error()), http.StatusBadRequest)
		http.Error(w, "Error decoding the request body.", http.StatusBadRequest)
		return
	}

	// CHECK if shorthand is ABC101 format
	if !hasThreeConsecutiveNumerics(courseShortHand) {
		logger.Error(fmt.Errorf("invalid Course shorthand"), http.StatusBadRequest)
		http.Error(w, "Invalid Course shorthand.", http.StatusBadRequest)
		return
	}
	update := bson.M{"$set": updateCourse}
	// Check if course exists
	filter = bson.D{{"shorthand", courseShortHand}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		logger.Error(fmt.Errorf("Error while updating the course: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error while updating the course.", http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated Successfuly")

	// Uncomment the follow line for debugging
	// logger.Info("UpdateCourse function completed.")
}

// DeleteCourse - deletes the course witht the given shorthand
func DeleteCourse(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("DeleteCourse function called.")

	// Extract the user username from the URL path
	path := r.URL.Path
	courseShortHand := strings.TrimPrefix(path, "/courses/")

	// CHECK if shorthand is ABC101 format
	if !hasThreeConsecutiveNumerics(courseShortHand) {
		logger.Error(fmt.Errorf("invalid Course shorthand"), http.StatusBadRequest)
		http.Error(w, "Invalid Course shorthand.", http.StatusBadRequest)
		return
	}

	// Check if course exists
	var result Course
	filter := bson.D{{"shorthand", courseShortHand}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		logger.Error(fmt.Errorf("Course %s doesn't exist", result.ShortHand), http.StatusInternalServerError)
		http.Error(w, fmt.Sprintf("Error: %s course doesn't exist", result.ShortHand), http.StatusInternalServerError)
		return
	}
	if err != nil {
		logger.Error(fmt.Errorf("Error while finding the course: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error while finding the course.", http.StatusInternalServerError)
		return
	}

	// Delete the course
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		logger.Error(fmt.Errorf("Error while Deleting the course: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error while Deleting the course.", http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted Successfuly")

	// Uncomment the follow line for debugging
	// logger.Info("DeleteCourse function completed.")
}

// checks for three consecutive digits
func hasThreeConsecutiveNumerics(input string) bool {
	re := regexp.MustCompile(`\d{3}`)
	return re.MatchString(input)
}
