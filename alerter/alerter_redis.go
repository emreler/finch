package alerter

import (
	"fmt"
	"regexp"
	"time"

	redis "gopkg.in/redis.v4"

	"github.com/emreler/finch/config"
	"github.com/emreler/finch/logger"
)

// RedisAlerter is the struct for alerting on event times
type RedisAlerter struct {
	client      *redis.Client
	alertChan   *chan string
	redisConfig *config.RedisConfig
	log         logger.InfoErrorLogger
}

// NewAlerter creates and returns new Alerter instance
func NewAlerter(config config.RedisConfig, c *chan string, l logger.InfoErrorLogger) *RedisAlerter {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Pwd,
		DB:       config.DB,
	})

	return &RedisAlerter{client: client, alertChan: c, redisConfig: &config, log: l}
}

// AddAlert method adds new alert to specified date
func (r *RedisAlerter) AddAlert(alertID string, alertDelay time.Duration) {
	r.client.Set(alertID, "1", 0)
	r.client.Expire(alertID, alertDelay)
}

// RemoveAlert removes alert
func (r *RedisAlerter) RemoveAlert(alertID string) {
	r.client.Del(alertID)
}

// StartListening starts to listen from Redis for alerts
func (r *RedisAlerter) StartListening() {
	go func() {
		for {
			// move from pending alerts queue to processing alerts queue
			msg := r.client.BRPopLPush(r.redisConfig.PendingAlertsKey, r.redisConfig.ProcessingAlertsKey, 0)
			alertID := string(msg.Val())

			// only send to channel if it looks like a mongo id and discard otherwise
			if match, _ := regexp.Match(`(?i)^[a-f\d]{24}$`, []byte(alertID)); match {
				*r.alertChan <- alertID
			} else {
				r.log.Error(fmt.Errorf("%s is not valid mongo id", alertID))
			}
		}
	}()
}

// RemoveProcessedAlert removes alerts from "currently processing alerts" queue
func (r *RedisAlerter) RemoveProcessedAlert(alertID string) {
	r.client.LRem(r.redisConfig.ProcessingAlertsKey, 0, alertID)
}

// AddAlertToQueue adds alerts to "process alerts" queue
func (r *RedisAlerter) AddAlertToQueue(alertID string) {
	r.client.LPush(r.redisConfig.PendingAlertsKey, alertID)
}
