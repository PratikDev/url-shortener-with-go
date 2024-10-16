package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pratikdev/url-shortner-with-go/controller/mongo_connection"
	"github.com/pratikdev/url-shortner-with-go/customErrors"
	"github.com/pratikdev/url-shortner-with-go/models"
	"github.com/pratikdev/url-shortner-with-go/token"
)

// HealthCheck shows the status of the server
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"ok"}`))
}

// HomeHandler is the handler for the home route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message":"Welcome to short-url"}`))
}

// LoginHandler is the handler for the login route
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginDetails models.LoginDetails
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		err = &customErrors.CustomError{Code: http.StatusBadRequest, Message: "Invalid request body"}
		customErrors.SendErrorResponse(w, err)
		return
	}

	// log the user in
	user, err := mongo_connection.LoginUser(loginDetails)
	if err != nil {
		customErrors.SendErrorResponse(w, err)
		return
	}

	// get token from the token module using username
	tokenResponse, err := token.GetToken(user.Username)
	if err != nil {
		customErrors.SendErrorResponse(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenResponse.Value,
		Expires: tokenResponse.ExpirationTime,
	})

	w.Write([]byte(`{"message":"Login Successful"}`))
}

// RegisterHandler is the handler for the register route
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte(`{"message":"Register Successful"}`))
}

// GetURL is the handler for the route that gets the URL from the id
func GetURL(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	url, err := mongo_connection.GetURLFromId(id)
	if err != nil {
		customErrors.SendErrorResponse(w, err)
		return
	}

	w.Write([]byte(`{"url":"` + url + `"}`))
}
