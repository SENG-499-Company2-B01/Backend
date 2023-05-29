package Courses

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/mongo"
)

type CourseModule struct{}

type Course struct {
	ShortHand     string   `json:"shorth"`
	Name          string   `json:"name"`
	Equipements   []string `json:"equipment"`
	Prerequisites []string `json:"prerequisites"`
}

func (cm CourseModule) CreateCourse(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	var newCourse Course
	err := json.NewDecoder(r.Body).Decode(&newCourse)
	if err != nil {
		http.Error(w, "Error while decoding "+err.Error(), http.StatusBadRequest)
		return
	}

	if !hasThreeConsecutiveNumerics(newCourse.ShortHand) {
		http.Error(w, "Invalid Course shorthand"+err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newCourse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newCourse)
}

// checks for three consecutive digits
func hasThreeConsecutiveNumerics(input string) bool {
	re := regexp.MustCompile(`\d{3}`)
	return re.MatchString(input)
}
