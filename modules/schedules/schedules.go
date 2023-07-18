package schedules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SENG-499-Company2-B01/Backend/logger"
	"github.com/SENG-499-Company2-B01/Backend/modules/classrooms"
	"github.com/SENG-499-Company2-B01/Backend/modules/courses"
	"github.com/SENG-499-Company2-B01/Backend/modules/users"
)

// Schedule represents a schedule entity
type Schedule struct {
	Year  string `json:"year"`
	Terms []Term `json:"terms"`
}

type Term struct {
	Term    string           `json:"term"`
	Courses []CourseOffering `json:"courses"`
}

type CourseOffering struct {
	Course   string  `json:"course"`
	Sections []Class `json:"sections"`
}

type Class struct {
	Num           string   `json:"num"`
	Building      string   `json:"building"`
	Professor     string   `json:"professor"`
	Days          []string `json:"days"`
	NumSeats      int      `json:"num_seats"`
	NumRegistered int      `json:"num_registered"`
	StartTime     string   `json:"start_time"`
	EndTime       string   `json:"end_time"`
}

type Algs2_Request struct {
	Year    string           `json:"year"`
	Term    string           `json:"term"`
	Courses []courses.Course `json:"courses"`
}

type Algs1_Request struct {
	Year       string                  `json:"year"`
	Term       string                  `json:"term"`
	Users      []users.User            `json:"users"`
	Courses    []CoursesWithCapacities `json:"courses"`
	Classrooms []classrooms.Classroom  `json:"classrooms"`
}

type Estimate struct {
	Course   string `json:"course"`
	Estimate int    `json:"estimate"`
}

type CoursesWithCapacities struct {
	ShortHand     string     `json:"shorthand" bson:"shorthand"`
	Name          string     `json:"name" bson:"name"`
	Prerequisites [][]string `json:"prerequisites" bson:"prerequisites"`
	CoRequisites  []string   `json:"corequisites" bson:"corequisites"`
	TermsOffered  []string   `json:"terms_offered" bson:"terms_offered"`
	Capacity      int        `json:"capacity" bson:"capacity"`
}

type Capacity struct {
	Estimate []Estimate `json:"estimate"`
}

// GenerateSchedule - Generates a new schedule
// TODO: Still needs to be implemented once algo team sets up their REST API
func GenerateSchedule(w http.ResponseWriter, r *http.Request, draft_schedules *mongo.Collection, users_coll *mongo.Collection, courses_coll *mongo.Collection, classrooms_coll *mongo.Collection, algs1_api string, algs2_api string) {
	logger.Info("GenerateSchedule function called.")

	// Extract the year and term values from the URL path
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 {
		// If there is an error parsing the url path,
		// log the error and return a bad request response
		logger.Error(fmt.Errorf("invalid URL path"), http.StatusBadRequest)
		http.Error(w, "Invalid URL path.", http.StatusBadRequest)
		return
	}

	// Extract the year and term from path
	year := path[2]
	term := path[3]

	// Check if passed term is valid
	if strings.ToLower(term) != "fall" && strings.ToLower(term) != "spring" && strings.ToLower(term) != "summer" {
		logger.Error(fmt.Errorf("invalid term for generating schedule"), http.StatusBadRequest)
		http.Error(w, "Invalid Term for Generating Schedule", http.StatusBadRequest)
		return
	}

	var courses_list []courses.Course

	// Retrieve all documents from the MongoDB collection
	cursor1, err := courses_coll.Find(context.TODO(), bson.M{"term": strings.ToLower(term)})
	if err != nil {
		logger.Error(fmt.Errorf("Error retrieving users: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving users.", http.StatusInternalServerError)
		return
	}
	defer cursor1.Close(context.TODO())

	// Iterate through the cursor and decode each document into a User struct
	for cursor1.Next(context.TODO()) {
		var course courses.Course
		err := cursor1.Decode(&course)
		if err != nil {
			logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
			return
		}
		courses_list = append(courses_list, course)
	}

	// Check for any errors during cursor iteration
	if err := cursor1.Err(); err != nil {
		logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
		return
	}

	// Create Algs 2 Request
	var new_algs2_request Algs2_Request
	new_algs2_request.Year = year
	new_algs2_request.Term = strings.ToLower(term)
	new_algs2_request.Courses = courses_list

	algs2RequestBody, _ := json.Marshal(new_algs2_request)
	algs2Payload := []byte(algs2RequestBody)
	algs2Req, _ := http.Post(algs2_api, "application/json", bytes.NewBuffer(algs2Payload))

	// Check the response status code and populate the capacity array
	var capacity []Capacity
	if algs2Req.StatusCode == http.StatusOK { // Response status is 200 (OK)
		// Parse the response body into the capacity variable
		decoder := json.NewDecoder(algs2Req.Body)
		err := decoder.Decode(&capacity)

		if err != nil {
			// Handle error in parsing response body
			// logger.Error(fmt.Errorf("Error trying to parse response body: "+err.Error()), http.StatusInternalServerError)
			// http.Error(w, "Error trying to parse response body.", http.StatusInternalServerError)

			// Construct an empty capacity array
			capacity = append(capacity, Capacity{})
		}
	} else { // Response status is not 200 (OK)
		// Construct an empty capacity array
		// create a random number between 80 and 100 for each course.
		capacity = append(capacity, Capacity{})
	}

	// Loop through the estimate object to create the CoursesWithCapacities object.
	for i := 0; i < len(capacity[0].Estimate); i++ {
		// Loop through the courses list to find the course with the same shorthand.
	}

	var users_list []users.User
	var classrooms_list []classrooms.Classroom

	cursor2, err := users_coll.Find(context.TODO(), bson.M{})
	if err != nil {
		logger.Error(fmt.Errorf("Error retrieving users: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving users.", http.StatusInternalServerError)
		return
	}
	defer cursor2.Close(context.TODO())

	// Iterate through the cursor and decode each document into a User struct
	for cursor2.Next(context.TODO()) {
		var user users.User
		err := cursor2.Decode(&user)
		if err != nil {
			logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
			return
		}
		users_list = append(users_list, user)
	}

	// Check for any errors during cursor iteration
	if err := cursor2.Err(); err != nil {
		logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
		return
	}

	// Retrieve all documents from the MongoDB collection
	cursor3, err := classrooms_coll.Find(context.TODO(), bson.M{})
	if err != nil {
		logger.Error(fmt.Errorf("Error retrieving classrooms: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving classrooms.", http.StatusInternalServerError)
		return
	}
	defer cursor3.Close(context.TODO())

	// Iterate through the cursor and decode each document into a User struct
	for cursor3.Next(context.TODO()) {
		var classroom classrooms.Classroom
		err := cursor3.Decode(&classroom)
		if err != nil {
			logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
			return
		}
		classrooms_list = append(classrooms_list, classroom)
	}

	// Check for any errors during cursor iteration
	if err := cursor3.Err(); err != nil {
		logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
		return
	}

	// Create Algs 1 Request
	var new_algs1_request Algs1_Request
	new_algs1_request.Year = year
	new_algs1_request.Term = term
	new_algs1_request.Users = users_list
	// new_algs1_request.Courses = courses_list
	new_algs1_request.Classrooms = classrooms_list

	algs1RequestBody, _ := json.Marshal(new_algs1_request)
	algs1Payload := []byte(algs1RequestBody)
	algs1Req, _ := http.Post(algs1_api, "application/json", bytes.NewBuffer(algs1Payload))

	var new_schedule Schedule
	final_err := json.NewDecoder(algs1Req.Body).Decode(&new_schedule)

	if final_err != nil {
		logger.Error(fmt.Errorf("Error parsing generated schedule: "+final_err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error parsing generated schedule.", http.StatusInternalServerError)
		return
	}

	// Store the schedule in the MongoDB collection
	_, insertErr := draft_schedules.InsertOne(context.TODO(), new_schedule)
	if insertErr != nil {
		logger.Error(fmt.Errorf("Error inserting schedule into collection: "+insertErr.Error()), http.StatusInternalServerError)
		http.Error(w, "Error inserting schedule into collection.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful schedule creation
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(new_schedule)
	// fmt.Fprintf(w, "New Scheduled generated successfully")

	// Uncomment the follow line for debugging
	// logger.Info("GenerateSchedule function completed.")
}

// ApproveSchedule - removes schedule in draft collection and adds it to previous_schedules collection, approving it.
func ApproveSchedule(w http.ResponseWriter, r *http.Request, draftsCollection *mongo.Collection, previousSchedulesCollection *mongo.Collection) {
	logger.Info("ApproveSchedule function called.")

	// Find the schedule in the "draft_schedule" collection
	var foundSchedule Schedule
	err := draftsCollection.FindOne(context.TODO(), bson.M{}).Decode(&foundSchedule)
	if err != nil {
		logger.Error(fmt.Errorf("failed to find schedule in drafts collection"), http.StatusInternalServerError)
		http.Error(w, "Failed to find schedule.", http.StatusInternalServerError)
		return
	}

	// Insert the found schedule into the "previous_schedules" collection
	_, err = previousSchedulesCollection.InsertOne(context.TODO(), foundSchedule)
	if err != nil {
		logger.Error(fmt.Errorf("failed to insert schedule into previous_schedules collection"), http.StatusInternalServerError)
		http.Error(w, "Failed to insert schedule.", http.StatusInternalServerError)
		return
	}

	// Delete the schedule from the "draft_schedule" collection
	_, err = draftsCollection.DeleteOne(context.TODO(), bson.M{})
	if err != nil {
		logger.Error(fmt.Errorf("failed to delete from drafts collection"), http.StatusInternalServerError)
		http.Error(w, "Failed to delete schedule.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful schedule creation
	w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Schedule has been approved")

	// Uncomment the follow line for debugging
	// logger.Info("ApproveSchedule function completed.")
}

// GetSchedules retrieves all schedules from the MongoDB collection
func GetSchedules(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetSchedules function called.")

	// Define an empty slice to store the schedules
	var schedules []Schedule

	// Retrieve all documents from the MongoDB collection
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		// If there is an error retrieving schedules,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error retrieving schedules: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error retrieving schedules.", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and decode each document into a Schedule struct
	for cursor.Next(context.TODO()) {
		var schedule Schedule
		err := cursor.Decode(&schedule)
		if err != nil {
			// If there is an error decoding a schedule document,
			// log the error and return an internal server error response
			logger.Error(fmt.Errorf("Error decoding schedule: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error decoding schedule.", http.StatusInternalServerError)
			return
		}
		schedules = append(schedules, schedule)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		// If there is an error iterating through the cursor,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error iterating cursor: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error iterating cursor.", http.StatusInternalServerError)
		return
	}

	// Send a response with the retrieved schedules
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedules)

	// Uncomment the follow line for debugging
	// logger.Info("GetSchedules function completed.")
}

// GetSchedule retrieves a schedule by year
func GetSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("GetSchedule function called.")

	// Extract the schedule year from the URL path
	path := r.URL.Path
	year := strings.TrimPrefix(path, "/schedules/")
	year = strings.TrimSpace(year)

	// Retrieve the schedule from the MongoDB collection
	filter := bson.M{"year": year}
	var schedule Schedule
	err := collection.FindOne(context.TODO(), filter).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If the schedule is not found,
			// log the error and return a not found response
			logger.Error(fmt.Errorf("Schedule not found: "+err.Error()), http.StatusNotFound)
			http.Error(w, "Schedule not found.", http.StatusNotFound)
		} else {
			// If there is an error retrieving the schedule,
			// log the error and return an internal server error response
			logger.Error(fmt.Errorf("Error getting schedule: "+err.Error()), http.StatusInternalServerError)
			http.Error(w, "Error getting schedule.", http.StatusInternalServerError)
		}
		return
	}

	// Send a response with the retrieved schedule
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)

	// Uncomment the follow line for debugging
	// logger.Info("GetSchedule function completed.")
}

// UpdateSchedule handles updating an existing schedule
func UpdateSchedule(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("UpdateSchedule function called.")

	// Extract the schedule year from the URL path
	path := r.URL.Path
	year := strings.TrimPrefix(path, "/schedules/")
	year = strings.TrimSpace(year)

	// Check if the year exists in the collection
	filter := bson.M{"year": year}
	exists, err := scheduleExists(filter, collection)
	if err != nil {
		// If there is an error querying the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error querying collection: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error querying collection.", http.StatusInternalServerError)
		return
	}
	if !exists {
		// If the year doesn't exist,
		// return a not found response
		logger.Error(fmt.Errorf("schedule does not exist"), http.StatusInternalServerError)
		http.Error(w, "Schedule does not exist.", http.StatusInternalServerError)
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

	// Construct the update query
	update := bson.M{"$set": requestBody}

	// Update the schedule in the MongoDB collection
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// If there is an error updating the schedule in the collection,
		// log the error and return an internal server error response
		logger.Error(fmt.Errorf("Error updating schedule: "+err.Error()), http.StatusInternalServerError)
		http.Error(w, "Error updating schedule.", http.StatusInternalServerError)
		return
	}

	// Send a response indicating successful schedule update
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Schedule updated successfully")

	// Uncomment the follow line for debugging
	// logger.Info("UpdateSchedule function completed.")
}

// scheduleExists checks if a document exists in the collection based on a filter
func scheduleExists(filter bson.M, collection *mongo.Collection) (bool, error) {
	count, err := collection.CountDocuments(context.TODO(), filter, nil)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
