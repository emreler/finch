package main

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/emreler/finch/config"
	"gitlab.com/emreler/finch/handler"
)

const prefix = "/v1"

func main() {
	config := config.NewConfig()

	storage := NewStorage(config.Mongo)

	c := make(chan string)
	alerter := NewAlerter(config.Redis, &c)

	handlers := InitHandlers(storage, alerter)

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
