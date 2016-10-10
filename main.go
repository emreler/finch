package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// MongoConfig has config values for MongoDB
type MongoConfig string

// RedisConfig has config values for Redis
type RedisConfig struct {
	Addr string `json:"addr"`
	Pwd  string `json:"pwd"`
}

// Config struct defines the config structure
type Config struct {
	Mongo MongoConfig `json:"mongo"`
	Redis RedisConfig `json:"redis"`
}

func main() {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("config.json file not found")
	}

	config := new(Config)
	json.Unmarshal(file, &config)

	storage := NewStorage(config.Mongo)
	c := make(chan string)
	alerter := NewAlerter(config.Redis, &c)
	handlers := InitHandlers(storage, alerter)

	alerter.StartListening()

	mux := http.NewServeMux()

	mux.Handle("/new-alert", http.HandlerFunc(handlers.NewAlert))

	go func() {
		for {
			alertID := <-c
			handlers.ProcessAlert(alertID)
		}
	}()

	log.Println("Starting server")
	fmt.Println(http.ListenAndServe(":8081", mux))

}
