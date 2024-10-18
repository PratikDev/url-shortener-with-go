package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/pratikdev/url-shortner-with-go/customErrors"
	"github.com/pratikdev/url-shortner-with-go/token"
)

type Middleware func(http.Handler) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := range middlewares {
			next = middlewares[len(middlewares)-1-i](next)
		}
		return next
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader is a method that implements the http.ResponseWriter interface
func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// middleware for logging requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(rw, r)
		log.Printf("%d %s %s %s", rw.statusCode, r.Method, r.URL.Path, time.Since(start))
	})
}

// middleware for recovering from panics
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"message": "` + http.StatusText(http.StatusInternalServerError) + `"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// middleware for setting headers for CORS
func SetHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// middleware for checking if the request is authenticated (takes a request handler as an argument)
func Auth(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// check if the request is authenticated
		cookieToken, err := r.Cookie("token")
		if err != nil {
			customErrors.SendErrorResponse(w, &customErrors.CustomError{Code: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}

		// refresh the token
		tokenResponse, err := token.RefreshToken(cookieToken.Value)
		if err != nil {
			customErrors.SendErrorResponse(w, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenResponse.Value,
			Expires: tokenResponse.ExpirationTime,
		})

		validationResponse, err := token.ValidateToken(tokenResponse.Value)
		if err != nil || !validationResponse.IsValid {
			customErrors.SendErrorResponse(w, err)
			return
		}

		// if the request is authenticated, set the userId in the request context
		r.Header.Set("userId", validationResponse.UserID)

		handler(w, r)
	}
}
