package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pratikdev/url-shortner-with-go/controller"
	"github.com/pratikdev/url-shortner-with-go/middleware"
)

const PORT = 5000

func main() {
	router := http.NewServeMux()

	// Public routes
	router.HandleFunc("GET /", controller.HomeHandler)
	router.HandleFunc("POST /login", controller.LoginHandler)
	router.HandleFunc("POST /register", controller.RegisterHandler)
	router.HandleFunc("GET /health", controller.HealthCheck)

	// Authenticated routes
	router.HandleFunc("GET /{id}", middleware.Auth(controller.GetURL))

	middlewareStack := middleware.Chain(middleware.Logging, middleware.Recover, middleware.SetHeaders)

	portString := ":" + strconv.Itoa(PORT)
	fmt.Printf("Server running in http://localhost:%d\n", PORT)
	log.Fatal(http.ListenAndServe(portString, middlewareStack(router)))
}
