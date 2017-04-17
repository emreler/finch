package storage

import "time"

// Alerter .
type Alerter interface {
	AddAlert(alertID string, alertDelay time.Duration)
	RemoveAlert(alertID string)
	StartListening()
	RemoveProcessedAlert(alertID string)
	AddAlertToQueue(alertID string)
}
