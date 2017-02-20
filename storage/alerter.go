package storage

import (
	"log"
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
		// before waiting for new alerts, handle the waiting ones in the list
		// that are created with ./persist-alerts/main.go
		alertsMap, _ := a.client.HGetAll(a.redisConfig.PendingAlertsHashKey).Result()

		for alertID := range alertsMap {
			*a.c <- string(alertID)
		}

		pubsub, err := a.client.Subscribe(a.redisConfig.AlertsChannelKey)

		if err != nil {
			panic(err)
		}

		for {
			msg, err := pubsub.ReceiveMessage()

			if err != nil {
				log.Println(err)
				continue
			}

			log.Println(string(msg.Payload))
			*a.c <- string(msg.Payload)
		}
	}()
}
