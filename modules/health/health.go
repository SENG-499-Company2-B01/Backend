package health

import (
	"fmt"
	"net/http"
)

// CheckHealth returns a 200OK response for the deployed server to occasionally check and monitor itself
func CheckHealth(w http.ResponseWriter, r *http.Request) {
	// Send a response indicating successful user creation
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "")
}
