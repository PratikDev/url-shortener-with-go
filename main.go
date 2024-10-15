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

	router.HandleFunc("GET /", controller.HomeHandler)
	router.HandleFunc("POST /login", controller.LoginHandler)
	router.HandleFunc("POST /register", controller.RegisterHandler)
	router.HandleFunc("GET /health", controller.HealthCheck)
	router.HandleFunc("GET /{id}", controller.GetURL)

	middlewareStack := middleware.Chain(middleware.Logging)

	portString := ":" + strconv.Itoa(PORT)
	server := http.Server{
		Addr:    portString,
		Handler: middlewareStack(router),
	}

	fmt.Printf("Server running in http://localhost:%d\n", PORT)

	log.Fatal(server.ListenAndServe())
}
