package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/emreler/finch/auth"
	"github.com/emreler/finch/config"
	"github.com/emreler/finch/handlers"
	"github.com/emreler/finch/logger"
	"github.com/emreler/finch/storage"
)

const prefix = "/v1"

func main() {
	configPath := flag.String("config", "", "Path of config.json file")
	flag.Parse()

	config := config.NewConfig(*configPath)

	c := make(chan string)

	auth := auth.NewAuth(config.Secret)
	stg := storage.NewStorage(config.Mongo)
	alerter := storage.NewAlerter(config.Redis, &c)
	logger := logger.NewLogger(config.Logentries)

	hnd := handlers.NewHandlers(stg, alerter, logger, auth)

	alerter.StartListening()

	mux := http.NewServeMux()

	// serving homepage
	mux.Handle("/", http.FileServer(http.Dir("web")))

	// serving api
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
