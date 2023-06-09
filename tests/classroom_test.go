package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/SENG-499-Company2-B01/Backend/modules/classrooms"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsertClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}
	setupRoutes(router)

	nClassroom.Shorthand = "Test"
	nClassroom.Building = "Test Building"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "12"

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	req, _ := http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	response := executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}
	print(response.Code)
	filter := bson.M{"shorthand": nClassroom.Shorthand, "room_number": nClassroom.Room_number}
	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter)
	})
}

func TestGetClassrooms(t *testing.T) {
	// Create a classroom
	var n1Classroom = classrooms.Classroom{}
	var n2Classroom = classrooms.Classroom{}
	setupRoutes(router)

	n1Classroom.Shorthand = "Test"
	n1Classroom.Building = "Test Building"
	n1Classroom.Capacity = 100
	n1Classroom.Room_number = "12"

	n2Classroom.Shorthand = "Test2"
	n2Classroom.Building = "Test Building"
	n2Classroom.Capacity = 100
	n2Classroom.Room_number = "12"

	requestBody, _ := json.Marshal(n1Classroom)
	payload := []byte(requestBody)
	req, _ := http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	response := executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	requestBody, _ = json.Marshal(n2Classroom)
	payload = []byte(requestBody)
	req, _ = http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	response = executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	req, _ = http.NewRequest("GET", "/classrooms", nil)
	response = executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	filter1 := bson.M{"shorthand": n1Classroom.Shorthand, "room_number": n1Classroom.Room_number}
	filter2 := bson.M{"shorthand": n2Classroom.Shorthand, "room_number": n2Classroom.Room_number}

	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter1)
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter2)
	})
}

func TestGetClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}

	setupRoutes(router)
	nClassroom.Shorthand = "Test"
	nClassroom.Building = "Test Building"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "2"

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	ins, _ := http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	insert_response := executeRequest(ins)

	if insert_response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, insert_response.Code)
	}

	get, _ := http.NewRequest("GET", "/classrooms/Test/2", nil)
	get_response := executeRequest(get)

	if get_response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, get_response.Code)
	}

	filter1 := bson.M{"shorthand": nClassroom.Shorthand, "room_number": nClassroom.Room_number}

	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter1)
	})
}

func TestDeleteClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}

	setupRoutes(router)
	nClassroom.Shorthand = "Test"
	nClassroom.Building = "Test Building"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "4"

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	req, _ := http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	response := executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	req, _ = http.NewRequest("DELETE", "/classrooms/Test/4", bytes.NewBuffer(payload))
	response = executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	filter := bson.M{"shorthand": nClassroom.Shorthand, "room_number": nClassroom.Room_number}

	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter)
	})
}

func TestUpdateClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}

	setupRoutes(router)
	nClassroom.Shorthand = "Test"
	nClassroom.Building = "Test Building"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "4"

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	req, _ := http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	response := executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	var n2Classroom = classrooms.Classroom{}

	n2Classroom = nClassroom

	n2Classroom.Building = "Test Building updated"
	requestBody, _ = json.Marshal(n2Classroom)
	payload = []byte(requestBody)
	req, _ = http.NewRequest("PUT", "/classrooms/Test/4", bytes.NewBuffer(payload))
	response = executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	req, _ = http.NewRequest("GET", "/classrooms/Test/4", bytes.NewBuffer(payload))
	response = executeRequest(req)
	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	var getClassroom = classrooms.Classroom{}

	json.Unmarshal([]byte(response.Body.String()), &getClassroom)

	if getClassroom.Building != "Test Building updated" {
		t.Errorf("Expected response body to be %s. Got %s\n", "Test Building updated", getClassroom.Building)
	}

	filter := bson.M{"shorthand": nClassroom.Shorthand, "room_number": nClassroom.Room_number}
	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter)
	})
}
