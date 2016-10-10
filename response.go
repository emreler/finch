package main

import (
	"encoding/json"
	"net/http"
)

// Response is the general response struct
type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// SendSuccess sends Response with {status: "success"}
func SendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&Response{Status: "success", Data: data})
}
