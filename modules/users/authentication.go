package users

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SENG-499-Company2-B01/Backend/logger"
)

// SignIn: Does Sign In process, and returns jwt token and user role
func SignIn(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	logger.Info("Signin function called.")

	// Define an empty slice to store the users
	var signInReq User
	err := json.NewDecoder(r.Body).Decode(&signInReq)
	if err != nil {
		logger.Error(fmt.Errorf("Error while decoding User object: " + err.Error()))
		http.Error(w, "Error while decoding User object.", http.StatusBadRequest)
		return
	}

	var user User
	// Retrieve the user credentials from the MongoDB collection
	filter := bson.M{"email": signInReq.Email}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Error(fmt.Errorf("username or password incorrect"))
			http.Error(w, "Username or password incorrect.", http.StatusNotFound)
			return
		}

		logger.Error(fmt.Errorf("Error while searching for user" + err.Error()))
		http.Error(w, "Error while searching for user.", http.StatusNotFound)
		return
	}

	if signInReq.Email != user.Email || signInReq.Password != user.Password {
		logger.Error(fmt.Errorf("username or Password Incorrect"))
		http.Error(w, "Username or password incorrect.", http.StatusNotFound)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["expiry"] = time.Now().Add(48 * time.Hour)
	claims["authorized"] = true
	claims["email"] = signInReq.Email
	claims["isAdmin"] = signInReq.IsAdmin

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		logger.Error(fmt.Errorf("Error while making JWT token"+ err.Error()))
		http.Error(w, "Error while making JWT token.", http.StatusNotFound)
		return
	}

	// Send a response with the retrieved users
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"jwt": tokenString})
}
