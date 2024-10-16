package customErrors

import (
	"fmt"
	"net/http"
)

// CustomError defines a custom error type
type CustomError struct {
	Code    int
	Message string
}

// Satisfy the error interface for CustomError
func (e *CustomError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Message)
}

func SendErrorResponse(w http.ResponseWriter, err error) {
	// set default error message and code
	responseMessage := "Something went wrong"
	responseCode := http.StatusInternalServerError

	// change error message and code if custom error
	if customErr, ok := err.(*CustomError); ok {
		responseMessage = customErr.Message
		responseCode = customErr.Code
	}

	w.WriteHeader(responseCode)

	w.Write([]byte(`{"message":"` + responseMessage + `"}`))
}
