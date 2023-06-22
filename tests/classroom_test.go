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

	nClassroom.Shorthand = "Test"
	nClassroom.Building = "Test Building"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "12"
	nClassroom.Equipment = []string{"Test Equipment 1", "Test Equipment 2"}

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	res, _ := http.Post("http://10.9.0.2:8000/classrooms", "application/json", bytes.NewBuffer(payload))

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, res.StatusCode)
	}
	print(res.StatusCode)
	filter := bson.M{"shorthand": nClassroom.Shorthand, "room_number": nClassroom.Room_number}
	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter)
	})
}

func TestGetClassrooms(t *testing.T) {

	//Simple GET test
	res, _ := http.Get("http://10.9.0.2:8000/classrooms")

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, res.StatusCode)
	}
}

func TestGetClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}

	setupRoutes(router)
	nClassroom.Shorthand = "Test"
	nClassroom.Building = "Test Building"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "2"
	nClassroom.Equipment = []string{"Test Equipment 1", "Test Equipment 2"}

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	res, _ := http.Post("http://10.9.0.2:8000/classrooms", "application/json", bytes.NewBuffer(payload))

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, res.StatusCode)
	}

	res, _ = http.Get("http://10.9.0.2:8000/classrooms/Test/2")
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, res.StatusCode)
	}

	filter1 := bson.M{"shorthand": nClassroom.Shorthand, "room_number": nClassroom.Room_number}

	t.Cleanup(func() {
		client.Database("schedule_db").Collection("classrooms").DeleteOne(context.TODO(), filter1)
	})
}

func TestDeleteClassroom(t *testing.T) {
	// Create a classroom
	var nClassroom = classrooms.Classroom{}

	//setupRoutes(router)
	nClassroom.Shorthand = "Test"
	nClassroom.Building = "Test Building"
	nClassroom.Capacity = 100
	nClassroom.Room_number = "5"
	nClassroom.Equipment = []string{"Test Equipment 1", "Test Equipment 2"}

	// requestBody, _ := json.Marshal(nClassroom)
	// payload := []byte(requestBody)

	// res, _ := http.Post("http://10.9.0.2:8000/classrooms", "application/json", bytes.NewBuffer(payload))

	// if res.StatusCode != http.StatusOK {
	// 	t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, res.StatusCode)
	// }

	req, err := http.NewRequest(http.MethodDelete, "http:/10.9.0.2:8000/classrooms/TEST/4", nil)
	if err != nil {
		panic(err)
	}
	response := executeRequest(req)

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
	nClassroom.Equipment = []string{"Test Equipment 1", "Test Equipment 2"}

	requestBody, _ := json.Marshal(nClassroom)
	payload := []byte(requestBody)

	req, err := http.NewRequest("POST", "/classrooms", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
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
