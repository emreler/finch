package handlers

import (
	"encoding/json"
	"net/http"
)

// FinchHandler .
type FinchHandler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func (fn FinchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))

		return
	}
	data, err := fn(w, r)

	if err != nil {
		SendError(w, err.Error())
	} else {
		SendSuccess(w, data)
	}
}

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

// SendError sends Response with {status: "error"}
func SendError(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&Response{Status: "error", Data: data})
}
