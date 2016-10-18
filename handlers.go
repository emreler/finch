package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"gopkg.in/mgo.v2/bson"

	"gitlab.com/emreler/finch/auth"
	"gitlab.com/emreler/finch/channel"
	"gitlab.com/emreler/finch/logger"
	"gitlab.com/emreler/finch/models"
)

const (
	// TypeHTTP .
	TypeHTTP = "http"
)

// Handlers .
type Handlers struct {
	stg    *Storage
	alt    *Alerter
	logger *logger.Logger
	auth   *auth.Auth
}

// InitHandlers initializes handlers
func InitHandlers(stg *Storage, alt *Alerter, logger *logger.Logger, auth *auth.Auth) *Handlers {
	h := &Handlers{stg: stg, alt: alt, logger: logger, auth: auth}

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

// Validate validates request
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
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// CreateAlert creates new alert
func (h *Handlers) CreateAlert(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	authorization := r.Header.Get("Authorization")

	re := regexp.MustCompile("Bearer (.*)")
	match := re.FindStringSubmatch(authorization)

	if len(match) != 2 {
		return nil, fmt.Errorf("No token found in requst")
	}

	tokenString := match[1]

	userID, err := h.auth.ValidateToken(tokenString)

	if err != nil {
		return nil, err
	}

	if r.Method == "POST" {
		req := &CreateAlertRequest{}

		decoder := json.NewDecoder(r.Body)

		if err = decoder.Decode(&req); err != nil {
			h.logger.Error(err)
			return nil, fmt.Errorf("Invalid format")
		}

		if err = req.Validate(); err != nil {
			h.logger.Error(err)
			return nil, err
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
			User:        bson.ObjectIdHex(userID),
		}

		if req.RepeatEvery > 0 {
			alert.Schedule = &models.Schedule{RepeatEvery: req.RepeatEvery}
		}

		alertID := h.stg.CreateAlert(alert)
		seconds := int(alertDate.Sub(time.Now()).Seconds())

		h.logger.Info(fmt.Sprintf("%d seconds later", seconds))

		h.alt.AddAlert(alertID, alertDate)

		res := &CreateAlertResponse{alertDate.Format(time.RFC3339)}
		return res, nil
	} else if r.Method == "GET" {

	} else {
		return nil, fmt.Errorf("Invalid method: %s", r.Method)
	}

	return nil, nil
}

// ProcessAlert processes the alert
func (h *Handlers) ProcessAlert(alertID string) {
	h.logger.Info(fmt.Sprintf("Getting %s", alertID))
	alert := h.stg.GetAlert(alertID)

	if alert.Channel == TypeHTTP {
		httpChannel := &channel.HttpChannel{}
		err := httpChannel.Notify(alert)

		if err != nil {
			h.logger.Info(fmt.Sprintf("Error while notifying with HTTP channel. %s", err.Error()))
		}

		if alert.Schedule != nil && alert.Schedule.RepeatEvery > 0 {
			nextAlertDate := time.Now().Add(time.Duration(alert.Schedule.RepeatEvery) * time.Second)
			h.logger.Info(fmt.Sprintf("Scheduling next alert %d seconds later at %s", alert.Schedule.RepeatEvery, nextAlertDate))
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
			h.logger.Error(err)
			return nil, fmt.Errorf("Invalid format")
		}

		if req.Name == nil || req.Email == nil {
			return nil, fmt.Errorf("Missing fields")
		}

		user := &User{Name: *req.Name, Email: *req.Email}

		userID, err := h.stg.CreateUser(user)

		if err != nil {
			h.logger.Error(err)
			return nil, err
		}

		exp := time.Now().Add(24 * 365 * time.Hour)
		tokenString, err := h.auth.GenerateToken(userID, exp)

		if err != nil {
			h.logger.Error(err)
			return nil, err
		}

		res := &CreateUserResponse{Token: tokenString, Expires: exp.Unix()}
		return res, nil
	}

	return nil, fmt.Errorf("Invalid method: %s", r.Method)
}
