package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"gitlab.com/emreler/finch/channel"
	"gitlab.com/emreler/finch/models"
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
	Token       string `json:"token"`
	Name        string `json:"name"`
	Channel     string `json:"channel"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"contentType"`
	Data        string `json:"data"`
	AlertDate   string `json:"alertDate"`
	AlertAfter  int    `json:"alertAfter"`
	RepeatEvery int    `json:"repeatEvery"`
}

func (r *CreateAlertRequest) Validate() error {
	if r.Channel == TypeHTTP {
		if match, _ := regexp.Match("^(http://|https://)", []byte(r.URL)); !match {
			return fmt.Errorf("url must be present in the http(s)://domain.com format")
		}

		if match, _ := regexp.Match("^(http://|https://)(localhost|127.0.0.1|172.17)", []byte(r.URL)); match {
			return fmt.Errorf("url can't be a local pointing address")
		}
	}
	return nil
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

		if err := decoder.Decode(&req); err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Invalid format")
		}

		if err := req.Validate(); err != nil {
			log.Println(err)
			return nil, err
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

		alert := &models.Alert{
			Name:        req.Name,
			AlertDate:   alertDate,
			Channel:     req.Channel,
			URL:         req.URL,
			Method:      req.Method,
			ContentType: req.ContentType,
			Data:        req.Data,
			User:        *userID,
		}

		if req.RepeatEvery > 0 {
			alert.Schedule = &models.Schedule{RepeatEvery: req.RepeatEvery}
		}

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
		httpChannel := &channel.HttpChannel{}
		err := httpChannel.Notify(alert)

		if err != nil {
			log.Printf("Error while notifying with HTTP channel. %s", err.Error())
		}

		if alert.Schedule != nil && alert.Schedule.RepeatEvery > 0 {
			nextAlertDate := time.Now().Add(time.Duration(alert.Schedule.RepeatEvery) * time.Second)
			log.Printf("Scheduling next alert %d seconds later at %s", alert.Schedule.RepeatEvery, nextAlertDate)
			h.alt.AddAlert(alertID, nextAlertDate)
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
