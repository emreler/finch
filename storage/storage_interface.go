package storage

import "github.com/emreler/finch/models"

// Storage .
type Storage interface {
	CreateUser(*models.User) error
	CreateAlert(*models.Alert) error
	GetAlert(string) (*models.Alert, error)
	UpdateAlert(*models.Alert) error
	GetUserAlerts(string) ([]*models.Alert, error)
	GetAlertHistory(alertID string, limit int) (*[]models.ProcessAlert, error)
	CountProcessAlertLogs() (int, error)
	LogProcessAlert(alert *models.Alert, statusCode int) error
	LogCreateAlert(alert *models.Alert) error
	LogCreateUser(user *models.User) error
}
