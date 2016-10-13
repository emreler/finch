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

// CreateAlertResponse .
type CreateAlertResponse struct {
	AlertDate string `json:"alertDate"`
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
func (h *Handlers) CreateAlert(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		req := &CreateAlertRequest{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Invalid format")
		}

		userID, err := h.stg.CheckToken(req.Token)

		if err != nil {
			return nil, fmt.Errorf("Invalid token")
		}

		var alertDate time.Time

		if req.AlertAfter != 0 {
			alertDate = time.Now().Add(time.Duration(req.AlertAfter) * time.Second)
		} else if req.AlertDate != "" {
			alertDate, err = time.Parse(time.RFC3339, req.AlertDate)

			if err != nil {
				return nil, fmt.Errorf("Invalid date format")
			}
		}

		alert := &Alert{Name: req.Name, AlertDate: alertDate, Channel: req.Channel, URL: req.URL, Data: req.Data, User: *userID}

		alertID := h.stg.CreateAlert(alert)
		seconds := int(alertDate.Sub(time.Now()).Seconds())

		log.Printf("%d seconds later", seconds)

		h.alt.AddAlert(alertID, alertDate)

		res := &CreateAlertResponse{alertDate.Format(time.RFC3339)}
		return res, nil
	}

	return nil, fmt.Errorf("Invalid method: %s", r.Method)
}

// ProcessAlert processes the alert
func (h *Handlers) ProcessAlert(alertID string) {
	log.Printf("Getting %s\n", alertID)
	alert := h.stg.GetAlert(alertID)
	log.Printf("Got %+v", alert)

	if alert.Channel == TypeHTTP {
		h := new(channel.HttpChannel)
		err := h.Notify(alert.URL, alert.Data)

		if err != nil {
			log.Printf("Error while notifying with HTTP channel. %s", err.Error())
		}
	}
}

// CreateUser .
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == "POST" {
		req := &CreateUserRequest{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Invalid format")
		}

		if req.Name == nil || req.Email == nil {
			return nil, fmt.Errorf("Missing fields")
		}

		user := &User{Name: *req.Name, Email: *req.Email}

		token, err := h.stg.CreateUser(user)

		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf(err.Error())
		}

		res := &CreateUserResponse{Token: token}
		return res, nil
	}

	return nil, fmt.Errorf("Invalid method: %s", r.Method)
}
