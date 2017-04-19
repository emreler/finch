package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/emreler/finch/auth"
	"github.com/emreler/finch/config"
	"github.com/emreler/finch/logger"
	"github.com/emreler/finch/models"
	"github.com/emreler/finch/storage"
)

const (
	typeHTTP    = "http"
	methodGet   = "GET"
	methodPost  = "POST"
	methodPatch = "PATCH"
)

// Handlers .
type Handlers struct {
	stg            storage.Storage
	alt            storage.Alerter
	logger         logger.InfoErrorLogger
	auth           *auth.Auth
	counterChannel chan bool
	appConfig      *config.AppConfig
}

// NewHandlers initializes handlers
func NewHandlers(stg storage.Storage, alt storage.Alerter, lgr logger.InfoErrorLogger, auth *auth.Auth, counterChannel chan bool, config *config.AppConfig) *Handlers {
	return &Handlers{
		stg:            stg,
		alt:            alt,
		logger:         lgr,
		auth:           auth,
		counterChannel: counterChannel,
		appConfig:      config,
	}
}

// extractUserID returns userID from Authorization header.
func (h *Handlers) extractUserID(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	re := regexp.MustCompile("Bearer (.*)")
	match := re.FindStringSubmatch(authHeader)

	if len(match) != 2 {
		return "", fmt.Errorf("No token found in requst")
	}

	tokenString := match[1]

	userID, err := h.auth.ValidateToken(tokenString)

	if err != nil {
		return "", err
	}

	return userID, nil
}

// AlertDetail returns alert object or it's history
func (h *Handlers) AlertDetail(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID, err := h.extractUserID(r)

	if err != nil {
		return nil, err
	}

	historyPattern := regexp.MustCompile("/alerts/([0-9A-Fa-f]{24})/history$")
	detailPattern := regexp.MustCompile("/alerts/([0-9A-Fa-f]{24})$")

	if match := historyPattern.FindStringSubmatch(r.URL.Path); len(match) == 2 && r.Method == methodGet {
		alertID := match[1]

		alert, err := h.stg.GetAlert(alertID)

		if err != nil {
			return nil, err
		}

		if alert.User.Hex() != userID {
			return nil, fmt.Errorf("Unauthorized request")
		}

		alertHistory, _ := h.stg.GetAlertHistory(alertID, h.appConfig.AlertLogLimit)

		return alertHistory, nil

	} else if match = detailPattern.FindStringSubmatch(r.URL.Path); len(match) == 2 {
		alertID := match[1]

		alert, err := h.stg.GetAlert(alertID)

		if err != nil {
			return nil, err
		}

		if alert.User.Hex() != userID {
			return nil, fmt.Errorf("Unauthorized request")
		}

		if r.Method == methodGet {
			return alert, nil
		} else if r.Method == methodPatch {
			req := &models.UpdateAlertRequest{}

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
	} else {
		return nil, fmt.Errorf("Invalid URL")
	}

}

// Alerts creates new alert
func (h *Handlers) Alerts(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID, err := h.extractUserID(r)

	if err != nil {
		return nil, err
	}

	if r.Method == methodPost {
		req := &models.CreateAlertRequest{}

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

		seconds := int64(alertDate.Sub(time.Now()).Seconds())

		if seconds == 0 {
			seconds = 1
		}

		h.logger.Info(fmt.Sprintf("%d seconds later", seconds))

		h.alt.AddAlert(alert.ID.Hex(), time.Duration(seconds)*time.Second)

		res := &models.CreateAlertResponse{AlertDate: alertDate.Format(time.RFC3339)}
		return res, nil
	} else if r.Method == "GET" {
		alerts, err := h.stg.GetUserAlerts(userID)

		if err != nil {
			h.logger.Error(err)
			return nil, err
		}

		res := &models.GetAlertsResponse{Alerts: alerts, Count: len(alerts)}

		return res, nil
	} else {
		return nil, fmt.Errorf("Invalid method: %s", r.Method)
	}
}

// CreateUser creates a new user
func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == methodPost {
		req := &models.CreateUserRequest{}

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

		res := &models.CreateUserResponse{Token: tokenString, Expires: exp.Unix()}
		return res, nil
	}

	return nil, fmt.Errorf("Invalid method: %s", r.Method)
}
