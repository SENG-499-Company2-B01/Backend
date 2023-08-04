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

	nClassroom.Building = "Test1"
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
	filter := bson.M{"building": nClassroom.Building, "room": nClassroom.Room_number}
	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter)
	})
}

func TestGetClassrooms(t *testing.T) {
	// Create a classroom
	var n1Classroom = classrooms.Classroom{}
	var n2Classroom = classrooms.Classroom{}
	setupRoutes(router)

	n1Classroom.Building = "Test2"
	n1Classroom.Capacity = 100
	n1Classroom.Room_number = "12"

	n2Classroom.Building = "Test3"
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

	filter1 := bson.M{"building": n1Classroom.Building, "room": n1Classroom.Room_number}
	filter2 := bson.M{"shorthand": n2Classroom.Building, "room": n2Classroom.Room_number}

	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter1)
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter2)
	})
}

func TestGetClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}

	setupRoutes(router)

	nClassroom.Building = "Test4"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "2"

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	ins, _ := http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	insert_response := executeRequest(ins)

	if insert_response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, insert_response.Code)
	}

	get, _ := http.NewRequest("GET", "/classrooms/Test4/2", nil)
	get_response := executeRequest(get)

	if get_response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, get_response.Code)
	}

	filter1 := bson.M{"building": nClassroom.Building, "room": nClassroom.Room_number}

	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter1)
	})
}

func TestDeleteClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}

	setupRoutes(router)

	nClassroom.Building = "Test5"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "4"

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	req, _ := http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	response := executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	req, _ = http.NewRequest("DELETE", "/classrooms/Test5/4", bytes.NewBuffer(payload))
	response = executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	filter := bson.M{"building": nClassroom.Building, "room": nClassroom.Room_number}

	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter)
	})
}

func TestUpdateClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}

	setupRoutes(router)
	nClassroom.Building = "Test6"
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

	n2Classroom.Capacity = 98
	requestBody, _ = json.Marshal(n2Classroom)
	payload = []byte(requestBody)
	req, _ = http.NewRequest("PUT", "/classrooms/Test6/4", bytes.NewBuffer(payload))
	response = executeRequest(req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	req, _ = http.NewRequest("GET", "/classrooms/Test6/4", bytes.NewBuffer(payload))
	response = executeRequest(req)
	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	var getClassroom = classrooms.Classroom{}

	json.Unmarshal([]byte(response.Body.String()), &getClassroom)

	if getClassroom.Capacity != 98 {
		t.Errorf("Expected response body to be %d. Got %d\n", 98, getClassroom.Capacity)
	}

	filter := bson.M{"building": nClassroom.Building, "room": nClassroom.Room_number}
	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter)
	})
}
