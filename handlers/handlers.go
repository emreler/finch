package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/emreler/finch/auth"
	"github.com/emreler/finch/channel"
	"github.com/emreler/finch/logger"
	"github.com/emreler/finch/models"
	"github.com/emreler/finch/storage"
)

const (
	// TypeHTTP .
	TypeHTTP = "http"
)

// Handlers .
type Handlers struct {
	stg            *storage.Storage
	alt            *storage.Alerter
	logger         *logger.Logger
	auth           *auth.Auth
	counterChannel chan bool
}

// NewHandlers initializes handlers
func NewHandlers(stg *storage.Storage, alt *storage.Alerter, logger *logger.Logger, auth *auth.Auth, counterChannel chan bool) *Handlers {
	return &Handlers{stg: stg, alt: alt, logger: logger, auth: auth, counterChannel: counterChannel}
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
	RepeatCount int    `json:"repeatCount"`
}

// Validate validates request
func (r *CreateAlertRequest) Validate() error {
	if r.Channel == TypeHTTP {
		if match, _ := regexp.Match("^(http://|https://)", []byte(r.URL)); !match {
			return fmt.Errorf("url must be present in the http(s)://domain.com format")
		}

		if match, _ := regexp.Match("^(http://|https://)(localhost|127.0.0.1|0.0.0.0|172.17)", []byte(r.URL)); match {
			return fmt.Errorf("url can't be a local pointing address")
		}

		if r.AlertAfter == 0 || r.AlertAfter < 0 && r.AlertDate == "" {
			return fmt.Errorf("Either 'alertAfter' or 'alertDate' fields must be present")
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

func (r *CreateUserRequest) Validate() error {
	if r.Name == nil || r.Email == nil {
		return fmt.Errorf("Missing fields")
	}

	*r.Name = strings.TrimSpace(*r.Name)
	*r.Email = strings.TrimSpace(*r.Email)

	if *r.Name == "" || *r.Email == "" {
		return fmt.Errorf("'name' and 'email' can't be empty")
	}

	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+$`)
	if !re.MatchString(*r.Email) {
		return fmt.Errorf("Invalid 'email' value")
	}

	return nil
}

// CreateUserResponse .
type CreateUserResponse struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// GetAlertsResponse .
type GetAlertsResponse struct {
	Alerts []*models.Alert `json:"alerts"`
	Count  int             `json:"count"`
}

type UpdateAlertRequest struct {
	Enabled *bool `json:"enabled"`
}

func (h *Handlers) AlertDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

	re = regexp.MustCompile("/alerts/([0-9A-Fa-f]{24})$")
	match = re.FindStringSubmatch(r.URL.Path)

	if len(match) != 2 {
		return nil, fmt.Errorf("Invalid URL")
	}

	alertID := match[1]

	alert, err := h.stg.GetAlert(alertID)

	if err != nil {
		return nil, err
	}

	if alert.User.Hex() != userID {
		return nil, fmt.Errorf("Unauthorized request")
	}

	if r.Method == "GET" {
		return alert, nil
	} else if r.Method == "PATCH" {
		req := &UpdateAlertRequest{}

		decoder := json.NewDecoder(r.Body)

		if err = decoder.Decode(&req); err != nil {
			h.logger.Error(err)
			return nil, fmt.Errorf("Invalid format")
		}

		if req.Enabled == nil {
			return nil, fmt.Errorf("Invalid format")
		}

		alert.Enabled = *req.Enabled

		err := h.stg.UpdateAlert(alert)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	return nil, fmt.Errorf("Invalid method")
}

// Alerts creates new alert
func (h *Handlers) Alerts(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	authorization := r.Header.Get("Authorization")

	re := regexp.MustCompile("Bearer (.*)")
	match := re.FindStringSubmatch(authorization)

	if len(match) != 2 {
		return nil, fmt.Errorf("No token found in request")
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

		if req.AlertAfter > 0 {
			alertDate = time.Now().Add(time.Duration(req.AlertAfter) * time.Second)
		} else if req.AlertDate != "" {
			alertDate, err = time.Parse(time.RFC3339, req.AlertDate)

			if err != nil {
				return nil, fmt.Errorf("Invalid date format")
			}
		}

		alert := models.NewAlert()
		alert.Name = req.Name
		alert.AlertDate = alertDate
		alert.Channel = req.Channel
		alert.URL = req.URL
		alert.Method = req.Method
		alert.ContentType = req.ContentType
		alert.Data = req.Data
		alert.User = bson.ObjectIdHex(userID)

		alert.Schedule = &models.Schedule{}

		if req.RepeatEvery > 0 {
			alert.Schedule.RepeatEvery = req.RepeatEvery
		}

		if req.RepeatCount > 0 {
			alert.Schedule.RepeatCount = req.RepeatCount
		} else {
			alert.Schedule.RepeatCount = -1
		}

		err := h.stg.CreateAlert(alert)

		if err != nil {
			h.logger.Error(err)
			return nil, err
		}

		h.stg.LogCreateAlert(alert)

		seconds := int(alertDate.Sub(time.Now()).Seconds())

		h.logger.Info(fmt.Sprintf("%d seconds later", seconds))

		h.alt.AddAlert(alert.ID.Hex(), alertDate)

		res := &CreateAlertResponse{alertDate.Format(time.RFC3339)}
		return res, nil
	} else if r.Method == "GET" {
		alerts, err := h.stg.GetUserAlerts(userID)

		if err != nil {
			h.logger.Error(err)
			return nil, err
		}

		res := &GetAlertsResponse{Alerts: alerts, Count: len(alerts)}

		return res, nil
	} else {
		return nil, fmt.Errorf("Invalid method: %s", r.Method)
	}
}

// ProcessAlert processes the alert
func (h *Handlers) ProcessAlert(alertID string) error {
	alert, err := h.stg.GetAlert(alertID)

	if err != nil {
		h.logger.Error(err)
		return err
	}

	if alert.Enabled == true && (alert.Schedule.RepeatCount == -1 || alert.Schedule.RepeatCount > 0) {
		h.logger.Info(fmt.Sprintf("Processing %s", alertID))
		if alert.Channel == TypeHTTP {
			httpChannel := &channel.HTTPChannel{}
			statusCode, err := httpChannel.Notify(alert)

			if err != nil {
				h.logger.Info(fmt.Sprintf("Error while notifying with HTTP channel. %s", err.Error()))
			}

			h.stg.LogProcessAlert(alert, statusCode)

			h.counterChannel <- true
		}

		if alert.Schedule.RepeatCount > 0 {
			alert.Schedule.RepeatCount--

			h.stg.UpdateAlert(alert)
		}
	}

	if alert.Schedule != nil && alert.Schedule.RepeatEvery > 0 {
		nextAlertDate := time.Now().Add(time.Duration(alert.Schedule.RepeatEvery) * time.Second)
		h.logger.Info(fmt.Sprintf("Scheduling next alert %d seconds later at %s", alert.Schedule.RepeatEvery, nextAlertDate))
		h.alt.AddAlert(alertID, nextAlertDate)
	}

	return nil
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

		if err = req.Validate(); err != nil {
			return nil, err
		}

		user := &models.User{Name: *req.Name, Email: *req.Email}

		err = h.stg.CreateUser(user)

		if err != nil {
			h.logger.Error(err)
			return nil, err
		}

		h.stg.LogCreateUser(user)

		exp := time.Now().Add(24 * 365 * time.Hour)
		tokenString, err := h.auth.GenerateToken(user.ID.Hex(), exp)

		if err != nil {
			h.logger.Error(err)
			return nil, err
		}

		res := &CreateUserResponse{Token: tokenString, Expires: exp.Unix()}
		return res, nil
	}

	return nil, fmt.Errorf("Invalid method: %s", r.Method)
}
