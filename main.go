package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pratikdev/url-shortner-with-go/controller"
)

const PORT = 5000

func main() {
	http.HandleFunc("GET /", controller.HomeHandler)
	http.HandleFunc("GET /{id}", controller.GetURL)

	fmt.Printf("Server running in http://localhost:%d\n", PORT)

	portString := ":" + strconv.Itoa(PORT)
	log.Fatal(http.ListenAndServe(portString, nil))
}
