package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type MongoConfig string
type RedisConfig struct {
	Addr string `json:"addr"`
	Pwd  string `json:"pwd"`
}

type Config struct {
	Mongo MongoConfig `json:"mongo"`
	Redis RedisConfig `json:"redis"`
}

func main() {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("config.json file not found")
		os.Exit(1)
	}
	config := new(Config)
	json.Unmarshal(file, &config)

	storage := NewStorage(config.Mongo)
	alerter := InitAlerter(config.Redis)
	handlers := InitHandlers(storage, alerter)

	c := make(chan string)
	alerter.Start(c)

	mux := http.NewServeMux()

	mux.Handle("/new", http.HandlerFunc(handlers.NewAlert))

	go func() {
		for {
			alertID := <-c
			handlers.ProcessAlert(alertID)
		}
	}()

	log.Println("Starting server")
	http.ListenAndServe(":8080", mux)

}
