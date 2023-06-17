package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SENG-499-Company2-B01/Backend/logger"
	"github.com/SENG-499-Company2-B01/Backend/modules/helper"
	"github.com/SENG-499-Company2-B01/Backend/modules/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Middleware function, which will be called for each request
func Users_API_Access_Control(next http.Handler, collection *mongo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ignore if this is not a call to users
		if !strings.Contains(r.URL.Path, "/users") {
			// Middleware successful
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized - Error no token provided", http.StatusUnauthorized)
			logger.Error(fmt.Errorf("Unauthorized - Error no token provided"))
			return
		}
		token, err := helper.CleanJWT(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized - Error "+err.Error(), http.StatusUnauthorized)
			logger.Error(fmt.Errorf("Error while cleaning jwt " + err.Error()))
			return
		}

		ok, jwtInfo, err := helper.VerifyJWT(token)
		if err != nil || !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			if err != nil {
				logger.Error(fmt.Errorf("Error while verifying jwt " + err.Error()))
			}
			return
		}
		fmt.Println(token)

		// user must be admin for CRUD operation on users
		if (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE") && !jwtInfo.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			logger.Error(fmt.Errorf("Forbidden - CRUD operation requested by " + jwtInfo.Email))
			return
		}

		// if user is not admin then they can only get their
		if (r.Method == "GET") && !jwtInfo.IsAdmin {
			// Extract the user username from the URL path
			path := r.URL.Path
			username := strings.TrimPrefix(path, "/users/")
			username = strings.TrimSpace(username)

			// NOTE: check for get all user from non admin

			// Retrieve the user from the MongoDB collection
			var user users.User
			filter := bson.M{"username": username}
			err := collection.FindOne(context.TODO(), filter).Decode(&user)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// If the user is not found,
					// log the error and return a not found response
					logger.Error(fmt.Errorf("User not found"))
					http.Error(w, "Forbidden", http.StatusForbidden)
				} else {
					// If there is an error retrieving the user,
					// log the error and return an internal server error response
					logger.Error(fmt.Errorf("Error getting user"))
					http.Error(w, "Error getting user from DB", http.StatusInternalServerError)
				}
				return
			}

			if user.Email != jwtInfo.Email {
				logger.Error(fmt.Errorf("Forbidden, non admin user trying to access other users info"))
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

		}

		// Middleware successful
		next.ServeHTTP(w, r)
	})

}
