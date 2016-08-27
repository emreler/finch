package main

import (
	"log"
	"time"

	redis "gopkg.in/redis.v4"
)

type Alerter struct {
	client *redis.Client
}

func InitAlerter(config RedisConfig) *Alerter {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Pwd,
		DB:       0,
	})

	return &Alerter{client: client}
}

func (a *Alerter) AddAlert(alertID string, alertDate time.Time) {
	a.client.Set(alertID, "1", 0)
	a.client.ExpireAt(alertID, alertDate)
}

func (a *Alerter) Start(c chan string) {
	go func() {
		pubsub, err := a.client.Subscribe("__keyevent@0__:expired")

		if err != nil {
			panic(err)
		}

		for {
			msg, err := pubsub.ReceiveMessage()

			if err != nil {
				panic(err)
			}

			log.Println(string(msg.Payload))
			c <- string(msg.Payload)
		}
	}()
}
