package middleware

import (
	"log"
	"net/http"
	"time"
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
