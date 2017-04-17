package storage

import "github.com/emreler/finch/models"

// Storage .
type Storage interface {
	CreateUser(*models.User) error
	CreateAlert(*models.Alert) error
	GetAlert(string) (*models.Alert, error)
	UpdateAlert(*models.Alert) error
	GetUserAlerts(string) ([]*models.Alert, error)
	GetAlertHistory(string, int) ([]*models.ProcessAlert, error)
	CountProcessAlertLogs() (int, error)
	LogProcessAlert(*models.Alert, int) error
	LogCreateAlert(*models.Alert) error
	LogCreateUser(*models.User) error
}
