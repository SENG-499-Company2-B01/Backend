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

func setupCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
}

func fetch_email(path string, collection *mongo.Collection) string {

	username := strings.TrimPrefix(path, "/users/")
	username = strings.TrimSpace(username)

	filter := bson.M{"username": username}
	var user users.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// If the user is not found,
			// log the error and return a not found response
			logger.Error(fmt.Errorf("user not found"), http.StatusNotFound)
		} else {
			logger.Error(fmt.Errorf("error getting user"), http.StatusInternalServerError)
		}
		return ""
	}

	return user.Email

}

func valid_permissions(r *http.Request, isAdmin bool, jwt_email string, collection *mongo.Collection) bool {

	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
		if isAdmin {
			return true
		}
		// If we have a non-admin token, a user can only update themselves
		return jwt_email == fetch_email(r.URL.Path, collection)
	}

	return true
}

// Middleware function, which will be called for each request
func Users_API_Access_Control(next http.Handler, collection *mongo.Collection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    setupCORS(&w, r)
		// ignore if this is not a call to users or prev schedules, classrooms
		if !strings.Contains(r.URL.Path, "/users") && !strings.Contains(r.URL.Path, "/courses") && !strings.Contains(r.URL.Path, "/schedules/prev") && !strings.Contains(r.URL.Path, "/schedules") && !strings.Contains(r.URL.Path, "/classrooms") {
			// Middleware successful
			next.ServeHTTP(w, r)
			return
		}

		apikey := r.Header.Get("apikey")
		if apikey != "" {
			check := helper.VerifyAPIKey(apikey)
			if check {
				// Middleware successful
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			logger.Error(fmt.Errorf("Unauthorized"), http.StatusUnauthorized)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized - Error no token provided", http.StatusUnauthorized)
			logger.Error(fmt.Errorf("Unauthorized - Error no token provided"), http.StatusUnauthorized)
			return
		}
		token, err := helper.CleanJWT(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized - Error "+err.Error(), http.StatusUnauthorized)
			logger.Error(fmt.Errorf("Error while cleaning jwt "+err.Error()), http.StatusUnauthorized)
			return
		}

		ok, jwtInfo, err := helper.VerifyJWT(token)
		if err != nil || !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			if err != nil {
				logger.Error(fmt.Errorf("Error while verifying jwt "+err.Error()), http.StatusUnauthorized)
			}
			return
		}

		// Role based access for courses endpoints
		if strings.Contains(r.URL.Path, "/courses") {

			// Get is allowed with valid jwt
			if r.Method == "GET" {
				next.ServeHTTP(w, r)
				return
			}

			// CRUD for courses is only allowed for admins
			if (r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE") && !jwtInfo.IsAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				logger.Error(fmt.Errorf("Forbidden - CRUD operation on courses  requested by "+jwtInfo.Email), http.StatusForbidden)
				return
			}
		}

		// Role based access for users endpoints
		if strings.Contains(r.URL.Path, "/classrooms") {

			// get forbidden for jwt
			if r.Method == "GET" {
				next.ServeHTTP(w, r)
				return
			}

			// user must be admin for CRUD operation on classrooms
			if (r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE") && !jwtInfo.IsAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				logger.Error(fmt.Errorf("Forbidden - CRUD operation requested by "+jwtInfo.Email), http.StatusForbidden)
				return
			}

		}

		// Role based access for previous schedules endpoints
		if r.URL.Path == "/schedules/prev" {
			if !jwtInfo.IsAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				logger.Error(fmt.Errorf("Forbidden - /schedules/prev CRUD operation requested by non admin - "+jwtInfo.Email), http.StatusForbidden)
				return
			}
		}

		// Role based access for schedules endpoints
		if strings.Contains(r.URL.Path, "/schedules") && r.URL.Path != "/schedules/prev" {
			// get forbidden for jwt
			if r.Method == "GET" {
				next.ServeHTTP(w, r)
				return
			}

			// user must be admin for CRUD operation on schedules
			if (r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE") && !jwtInfo.IsAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				logger.Error(fmt.Errorf("Forbidden - CRUD operation requested for schedules by "+jwtInfo.Email), http.StatusForbidden)
				return
			}
		}

		// Role Based access for users endpoint
		if strings.Contains(r.URL.Path, "/users") {
			// user must be admin for CRUD operation on users
			if !valid_permissions(r, jwtInfo.IsAdmin, jwtInfo.Email, collection) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				logger.Error(fmt.Errorf("Forbidden - CRUD operation requested by "+jwtInfo.Email), http.StatusForbidden)
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
						logger.Error(fmt.Errorf("User not found"), http.StatusNotFound)
						http.Error(w, "Forbidden", http.StatusForbidden)
					} else {
						// If there is an error retrieving the user,
						// log the error and return an internal server error response
						logger.Error(fmt.Errorf("Error getting user"), http.StatusInternalServerError)
						http.Error(w, "Error getting user from DB", http.StatusInternalServerError)
					}
					return
				}

				if user.Email != jwtInfo.Email {
					logger.Error(fmt.Errorf("Forbidden, non admin user trying to access other users info"), http.StatusForbidden)
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

			}
		}
		// Middleware successful
		next.ServeHTTP(w, r)
	})

}
