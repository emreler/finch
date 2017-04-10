package models

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// TypeHTTP .
	TypeHTTP = "http"
)

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

// Validate validates CreateAlertRequest
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
	} else {
		return fmt.Errorf("Unsupported channel: %s", r.Channel)
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

// Validate validates CreateUserRequest
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
	Alerts []*Alert `json:"alerts"`
	Count  int      `json:"count"`
}

// UpdateAlertRequest .
type UpdateAlertRequest struct {
	Enabled *bool `json:"enabled"`
}
