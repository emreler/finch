package main

import (
	"flag"
	"log"

	"github.com/emreler/finch/config"

	redis "gopkg.in/redis.v4"
)

func main() {
	configPath := flag.String("config", "", "Path of config.json file")
	flag.Parse()

	config := config.NewConfig(*configPath)

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

	log.Println("Starting finch-persist-alerts")

	for {
		alert, err := pubsub.ReceiveMessage()

		if err != nil {
			log.Print(err)
			continue
		}

		alertID := alert.Payload

		log.Println(alertID)

		go func() {
			err := client.RPush(config.Redis.AlertsChannelKey, alertID).Err()
			if err != nil {
				log.Print(err)
			}
		}()
	}
}
