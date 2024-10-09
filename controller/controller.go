package controller

import (
	"encoding/json"
	"net/http"

	"github.com/pratikdev/url-shortner-with-go/controller/mongo_connection"
	"github.com/pratikdev/url-shortner-with-go/customErrors"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": "something here"})
}

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
