package storage

import (
	"time"

	redis "gopkg.in/redis.v4"

	"github.com/emreler/finch/config"
)

// Alerter is the struct for alerting on event times
type Alerter struct {
	client      *redis.Client
	c           *chan string
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

	return &Alerter{client: client, c: c, redisConfig: &config}
}

// AddAlert method adds new alert to specified date
func (a *Alerter) AddAlert(alertID string, alertDate time.Time) {
	a.client.Set(alertID, "1", 0)
	a.client.ExpireAt(alertID, alertDate)
}

// RemoveAlert removes alert
func (a *Alerter) RemoveAlert(alertID string) {
	a.client.Del(alertID)
}

// StartListening starts to listen from Redis for alerts
func (a *Alerter) StartListening() {
	go func() {
		for {
			msg := a.client.BLPop(0, a.redisConfig.AlertsChannelKey)
			alertID := msg.Val()[1]

			*a.c <- string(alertID)
		}
	}()
}
