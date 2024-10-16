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
	router.HandleFunc("GET /health", controller.HealthCheck)
	router.HandleFunc("GET /{id}", controller.GetURL)
	router.HandleFunc("POST /login", controller.LoginHandler)
	router.HandleFunc("POST /register", controller.RegisterHandler)

	// Authenticated routes
	router.HandleFunc("GET /all", middleware.Auth(controller.GetAllURL))
	router.HandleFunc("POST /new", middleware.Auth(controller.NewURL))

	middlewareStack := middleware.Chain(middleware.Logging, middleware.Recover, middleware.SetHeaders)

	portString := ":" + strconv.Itoa(PORT)
	fmt.Printf("Server running in http://localhost:%d\n", PORT)
	log.Fatal(http.ListenAndServe(portString, middlewareStack(router)))
}
