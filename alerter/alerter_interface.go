package alerter

import "time"

// Alerter .
type Alerter interface {
	AddAlert(string, time.Duration)
	RemoveAlert(string)
	StartListening()
	RemoveProcessedAlert(string)
	AddAlertToQueue(string)
}
