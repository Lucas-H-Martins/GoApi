package main

import (
	"fmt"
	"log"
	"net/http"
	"goapi/routes"
)

func main() {
	router := routes.SetupRoutes()
	
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
} 