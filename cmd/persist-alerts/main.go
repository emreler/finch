package main

import (
	"flag"

	"github.com/emreler/finch/config"
	"github.com/emreler/finch/logger"

	redis "gopkg.in/redis.v4"
)

func main() {
	configPath := flag.String("config", "", "Path of config.json file")
	flag.Parse()

	appLogger := logger.NewLogger()

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

	appLogger.Info("Starting finch-persist-alerts")

	for {
		alert, err := pubsub.ReceiveMessage()

		if err != nil {
			appLogger.Error(err)
			continue
		}

		alertID := alert.Payload

		appLogger.Info(alertID)

		go func() {
			err := client.LPush(config.Redis.PendingAlertsKey, alertID).Err()
			if err != nil {
				appLogger.Error(err)
			}
		}()
	}
}
