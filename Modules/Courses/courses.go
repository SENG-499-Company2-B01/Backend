package courses

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Course struct {
	ShortHand     string   `json:"shorth" bson:"shorthand"`
	Name          string   `json:"name" bson:"name"`
	Equipements   []string `json:"equipment" bson:"equipements"`
	Prerequisites []string `json:"prerequisites" bson:"prerequisites"`
}

func CreateCourse(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	var newCourse Course
	err := json.NewDecoder(r.Body).Decode(&newCourse)
	if err != nil {
		http.Error(w, "Error while decoding "+err.Error(), http.StatusBadRequest)
		return
	}

	// CHECK if shorthand is ABC101 format
	if !hasThreeConsecutiveNumerics(newCourse.ShortHand) {
		http.Error(w, "Invalid Course shorthand", http.StatusBadRequest)
		return
	}

	// CHECK if course doesn't exist
	var result Course
	filter := bson.D{{"shorthand", newCourse.ShortHand}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil && result.ShortHand != "" {
		http.Error(w, fmt.Sprintf("Error: %s course already exists", result.ShortHand), http.StatusInternalServerError)
		return
	}
	if err != nil && err != mongo.ErrNoDocuments {
		http.Error(w, "Error while checking for duplicate document, " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the user into the MongoDB collection
	_, err = collection.InsertOne(context.TODO(), newCourse)
	if err != nil {
		http.Error(w, "Error while inserting course into DB" + err.Error(), http.StatusInternalServerError)
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
