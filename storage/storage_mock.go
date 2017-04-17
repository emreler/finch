package storage

import (
	"github.com/emreler/finch/models"
	"gopkg.in/mgo.v2/bson"
)

// MockStorage is the mock implementation of the MongoStorage struct
type MockStorage struct {
	users  map[string]*models.User
	alerts map[string]*models.Alert
}

// NewMockStorage initalizes and returns new MockStorage instance
func NewMockStorage() *MockStorage {
	return &MockStorage{
		users:  make(map[string]*models.User),
		alerts: make(map[string]*models.Alert),
	}
}

// CreateUser .
func (m *MockStorage) CreateUser(user *models.User) error {
	user.ID = bson.NewObjectId()
	m.users[user.ID.Hex()] = user

	return nil
}

// CreateAlert .
func (m *MockStorage) CreateAlert(alert *models.Alert) error {
	alert.ID = bson.NewObjectId()
	m.alerts[alert.ID.Hex()] = alert

	return nil
}

// GetAlert .
func (m *MockStorage) GetAlert(alertID string) (*models.Alert, error) {
	return m.alerts[alertID], nil
}

// UpdateAlert .
func (m *MockStorage) UpdateAlert(alert *models.Alert) error {
	m.alerts[alert.ID.Hex()] = alert
	return nil
}

// GetUserAlerts .
func (m *MockStorage) GetUserAlerts(userID string) ([]*models.Alert, error) {
	var res []*models.Alert
	for _, v := range m.alerts {
		if v.User.Hex() == userID {
			res = append(res, v)
		}
	}

	return res, nil
}

// GetAlertHistory .
func (m *MockStorage) GetAlertHistory(alertID string, limit int) ([]*models.ProcessAlert, error) {
	var res []*models.ProcessAlert
	return res, nil
}

// CountProcessAlertLogs .
func (m *MockStorage) CountProcessAlertLogs() (int, error) {
	return 17, nil
}

// LogProcessAlert .
func (m *MockStorage) LogProcessAlert(alert *models.Alert, statusCode int) error {
	return nil
}

// LogCreateAlert .
func (m *MockStorage) LogCreateAlert(alert *models.Alert) error {
	return nil
}

// LogCreateUser .
func (m *MockStorage) LogCreateUser(alert *models.User) error {
	return nil
}
