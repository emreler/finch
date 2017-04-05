package storage

import (
	"log"
	"regexp"
	"time"

	redis "gopkg.in/redis.v4"

	"github.com/emreler/finch/config"
)

// Alerter is the struct for alerting on event times
type Alerter struct {
	client      *redis.Client
	alertIDChan *chan string
	redisConfig *config.RedisConfig
}

// NewAlerter creates and returns new Alerter instance
func NewAlerter(config config.RedisConfig, c *chan string) *Alerter {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Pwd,
		DB:       0,
	})

	client.ConfigSet("notify-keyspace-events", "Ex")

	return &Alerter{client: client, alertIDChan: c, redisConfig: &config}
}

// AddAlert method adds new alert to specified date
func (a *Alerter) AddAlert(alertID string, alertDelay time.Duration) {
	a.client.Set(alertID, "1", 0)
	a.client.Expire(alertID, alertDelay)
}

// RemoveAlert removes alert
func (a *Alerter) RemoveAlert(alertID string) {
	a.client.Del(alertID)
}

// StartListening starts to listen from Redis for alerts
func (a *Alerter) StartListening() {
	go func() {
		for {
			// move from pending alerts queue to processing alerts queue
			msg := a.client.BRPopLPush(a.redisConfig.PendingAlertsKey, a.redisConfig.ProcessingAlertsKey, 0)
			alertID := string(msg.Val())

			// only send to channel if it looks like a mongo id and discard otherwise
			if match, _ := regexp.Match(`(?i)^[a-f\d]{24}$`, []byte(alertID)); match {
				*a.alertIDChan <- alertID
			} else {
				log.Printf("%s is not a valid mongo id", alertID)
			}
		}
	}()
}

// RemoveProcessedAlert removes alerts from "currently processing alerts" queue
func (a *Alerter) RemoveProcessedAlert(alertID string) {
	a.client.LRem(a.redisConfig.ProcessingAlertsKey, 0, alertID)
}

// AddAlertToQueue adds alerts to "process alerts" queue
func (a *Alerter) AddAlertToQueue(alertID string) {
	a.client.LPush(a.redisConfig.PendingAlertsKey, alertID)
}
