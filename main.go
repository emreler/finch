package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/emreler/finch/auth"
	"gitlab.com/emreler/finch/config"
	"gitlab.com/emreler/finch/handlers"
	"gitlab.com/emreler/finch/logger"
	"gitlab.com/emreler/finch/storage"
)

const prefix = "/v1"
const defaultConfigPath = "/etc/finch/config.json"

func main() {
	configPath := flag.String("config", defaultConfigPath, "Path of config.json file")
	flag.Parse()

	c := make(chan string)

	config := config.NewConfig(*configPath)

	auth := auth.NewAuth(config.Secret)
	stg := storage.NewStorage(config.Mongo)
	alerter := storage.NewAlerter(config.Redis, &c)
	logger := logger.NewLogger(config.Logentries)

	hnd := handlers.NewHandlers(stg, alerter, logger, auth)

	alerter.StartListening()

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("web")))
	mux.Handle(prefix+"/alerts/", handlers.FinchHandler(hnd.AlertDetail))
	mux.Handle(prefix+"/alerts", handlers.FinchHandler(hnd.Alerts))
	mux.Handle(prefix+"/users", handlers.FinchHandler(hnd.CreateUser))

	go func() {
		for {
			alertID := <-c
			hnd.ProcessAlert(alertID)
		}
	}()

	log.Println("Starting server")
	fmt.Println(http.ListenAndServe(":8081", mux))

}
