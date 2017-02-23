package storage

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/emreler/finch/models"
)

const eventsCollection = "events"
const typeProcessAlert = "process_alert"

// ProcessAlert creates new user
func (s *Storage) LogProcessAlert(alert *models.Alert, statusCode int) error {
	ses := s.GetDBSession()
	defer ses.Close()

	data := struct {
		Type       string
		Alert      bson.ObjectId
		StatusCode int
		CreatedAt  time.Time
	}{
		typeProcessAlert,
		alert.ID,
		statusCode,
		time.Now(),
	}

	err := ses.DB("finch").C(eventsCollection).Insert(data)

	if err != nil {
		return err
	}

	return nil
}
