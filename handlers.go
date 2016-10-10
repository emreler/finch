package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.com/emreler/finch/channel"
)

const (
	// TypeHTTP .
	TypeHTTP = "http"
)

// Handlers .
type Handlers struct {
	stg *Storage
	alt *Alerter
}

// InitHandlers initializes handlers
func InitHandlers(stg *Storage, alt *Alerter) *Handlers {
	h := &Handlers{stg: stg, alt: alt}

	return h
}

// CreateAlertRequest .
type CreateAlertRequest struct {
	Token      string `json:"token"`
	Name       string `json:"name"`
	Channel    string `json:"channel"`
	URL        string `json:"url"`
	Data       string `json:"data"`
	AlertDate  string `json:"alertDate"`
	AlertAfter int    `json:"alertAfter"`
}

// CreateUserRequest .
type CreateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// CreateUserResponse .
type CreateUserResponse struct {
	Token string `json:"token"`
}

// CreateAlert creates new alert
func (h *Handlers) CreateAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req := &CreateAlertRequest{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "invalid format")
			return
		}

		userID, err := h.stg.CheckToken(req.Token)

		if err != nil {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		var alertDate time.Time

		if req.AlertAfter != 0 {
			alertDate = time.Now().Add(time.Duration(req.AlertAfter) * time.Second)
		} else if req.AlertDate != "" {
			alertDate, err = time.Parse("2006-01-02T15:04:05Z", req.AlertDate)

			if err != nil {
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return
			}
		}

		alert := &Alert{Name: req.Name, AlertDate: alertDate, Channel: req.Channel, URL: req.URL, Data: req.Data, User: *userID}

		alertID := h.stg.CreateAlert(alert)
		seconds := int(alertDate.Sub(time.Now()).Seconds())
		log.Printf("%d seconds later", seconds)

		h.alt.AddAlert(alertID, alertDate)
	}
}

// ProcessAlert processes the alert
func (h *Handlers) ProcessAlert(alertID string) {
	log.Printf("Getting %s\n", alertID)
	alert := h.stg.GetAlert(alertID)
	log.Printf("Got %+v", alert)

	if alert.Channel == TypeHTTP {
		h := new(channel.HttpChannel)
		h.Notify(alert.URL, alert.Data)
	}
}

// CreateUser .
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req := &CreateUserRequest{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "invalid format")
			return
		}

		if req.Name == nil || req.Email == nil {
			fmt.Fprint(w, "invalid format")
			return
		}

		user := &User{Name: *req.Name, Email: *req.Email}

		token, err := h.stg.CreateUser(user)

		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
		}

		SendSuccess(w, &CreateUserResponse{Token: token})
	}
}
