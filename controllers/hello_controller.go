package controllers

import (
	"encoding/json"
	"net/http"
	"goapi/models"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	response := models.Response{
		Message: "Hello, World!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
} 