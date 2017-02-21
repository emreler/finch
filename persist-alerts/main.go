package main

import (
	"log"

	"github.com/emreler/finch/config"

	redis "gopkg.in/redis.v4"
)

func main() {
	config := config.NewConfig("../config.json")

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Pwd,
		DB:       0,
	})

	client.ConfigSet("notify-keyspace-events", "Ex")

	// redis' key expiration channel, this is enabled by the config line "notify-keyspace-events Ex"
	pubsub, err := client.Subscribe("__keyevent@0__:expired")

	if err != nil {
		panic(err)
	}

	for {
		msg, err := pubsub.ReceiveMessage()

		if err != nil {
			panic(err)
		}

		alertID := msg.Payload

		log.Printf("received %s", alertID)

		listeners, _ := client.Publish(config.Redis.AlertsChannelKey, alertID).Result()

		// if there is no subscribers to the channel, that means finch is not running
		// so we persist them to a list, which is going to be scanned on startup by finch

		log.Printf("delivered to %d listeners", listeners)

		if listeners == 0 {
			err := client.HSet(config.Redis.PendingAlertsHashKey, alertID, "1").Err()

			if err != nil {
				log.Println(err)
			}
		}
	}
}
