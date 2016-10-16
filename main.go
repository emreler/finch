package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/emreler/finch/config"
	"gitlab.com/emreler/finch/handler"
	"gitlab.com/emreler/finch/logger"
)

const prefix = "/v1"
const defaultConfigPath = "/etc/finch/config.json"

func main() {
	configPath := flag.String("config", defaultConfigPath, "Path of config.json file")
	flag.Parse()

	config := config.NewConfig(*configPath)

	storage := NewStorage(config.Mongo)

	c := make(chan string)
	alerter := NewAlerter(config.Redis, &c)

	logger := logger.NewLogger(config.Logentries)

	handlers := InitHandlers(storage, alerter, logger)

	alerter.StartListening()

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("web")))
	mux.Handle(prefix+"/alerts", handler.FinchHandler(handlers.CreateAlert))
	mux.Handle(prefix+"/users", handler.FinchHandler(handlers.CreateUser))

	go func() {
		for {
			alertID := <-c
			handlers.ProcessAlert(alertID)
		}
	}()

	log.Println("Starting server")
	fmt.Println(http.ListenAndServe(":8081", mux))

}
