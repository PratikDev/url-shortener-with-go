package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pratikdev/url-shortner-with-go/controller/mongo_connection"
	"github.com/pratikdev/url-shortner-with-go/customErrors"
	"github.com/pratikdev/url-shortner-with-go/models"
)

// HealthCheck shows the status of the server
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// HomeHandler is the handler for the home route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": "something here"})
}

// LoginHandler is the handler for the login route
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var loginDetails models.LoginDetails
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		err = &customErrors.CustomError{Code: http.StatusBadRequest, Message: "Invalid request body"}
		customErrors.SendErrorResponse(w, err)
		return
	}

	if _, err := mongo_connection.LoginUser(loginDetails); err != nil {
		customErrors.SendErrorResponse(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// RegisterHandler is the handler for the register route
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var registerDetails models.LoginDetails
	if err := json.NewDecoder(r.Body).Decode(&registerDetails); err != nil {
		err = &customErrors.CustomError{Code: http.StatusBadRequest, Message: "Invalid request body"}
		customErrors.SendErrorResponse(w, err)
		return
	}

	if err := mongo_connection.RegisterUser(registerDetails); err != nil {
		customErrors.SendErrorResponse(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Register successful"})
}

// GetURL is the handler for the route that gets the URL from the id
func GetURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := r.PathValue("id")
	url, err := mongo_connection.GetURLFromId(id)
	if err != nil {
		customErrors.SendErrorResponse(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
