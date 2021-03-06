package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/emreler/finch/alerter"
	"github.com/emreler/finch/auth"
	"github.com/emreler/finch/config"
	"github.com/emreler/finch/counter"
	"github.com/emreler/finch/errors"
	"github.com/emreler/finch/handlers"
	"github.com/emreler/finch/logger"
	"github.com/emreler/finch/storage"
	"github.com/gorilla/websocket"
)

const prefix = "/v1"

func main() {
	configPath := flag.String("config", "", "Path of config.json file")
	flag.Parse()

	config := config.NewConfig(*configPath)

	alertChannel := make(chan string)
	counterChannel := make(chan bool)

	appLogger := logger.NewLogger(os.Stderr)
	auth := auth.NewAuth(config.Secret)
	stg := storage.NewStorage(config.Mongo)
	alt := alerter.NewAlerter(config.Redis, &alertChannel, appLogger)
	hnd := handlers.NewHandlers(stg, alt, appLogger, auth, counterChannel, &config.App)

	processedAlertCount, err := stg.CountProcessAlertLogs()
	if err != nil {
		panic(err)
	}

	alt.StartListening()

	mux := http.NewServeMux()

	// serving homepage
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/" {
			// request for index.html, parse template with counter value
			t := template.New("index.html")
			t, err := t.ParseFiles("web/index.html")

			if err != nil {
				appLogger.Error(err)
			}

			vars := struct {
				Counter int
			}{
				processedAlertCount,
			}

			t.Execute(w, vars)
			return
		}

		http.FileServer(http.Dir("web")).ServeHTTP(w, r)
	})

	// serving api
	mux.Handle(prefix+"/alerts/", handlers.FinchHandler(hnd.AlertDetail))
	mux.Handle(prefix+"/alerts", handlers.FinchHandler(hnd.Alerts))
	mux.Handle(prefix+"/users", handlers.FinchHandler(hnd.CreateUser))

	hub := counter.NewHub()
	go hub.Run()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// new incoming ws connction
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true }, // allow connections from all origins
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			appLogger.Error(err)
			return
		}

		client := &counter.Client{Conn: conn, Send: make(chan []byte)}
		hub.Register <- client

		client.WaitMessages()
	})

	go func() {
		ticker := time.NewTicker(time.Millisecond * 500)
		lastSent := processedAlertCount
		for {
			select {
			case <-counterChannel:
				processedAlertCount++
			case <-ticker.C:
				if processedAlertCount > lastSent {
					// increment the counter on clients if necessary
					hub.Broadcast <- []byte(strconv.Itoa(processedAlertCount))
					lastSent = processedAlertCount
				}
			}
		}
	}()

	go func() {
		for {
			alertID := <-alertChannel

			go func(alertID string) {
				err := hnd.ProcessAlert(alertID)

				if err == nil {
					alt.RemoveProcessedAlert(alertID)
				} else if _, ok := err.(*errors.RetryProcessError); ok {
					appLogger.Info("retrying")
					appLogger.Error(err)
					alt.AddAlertToQueue(alertID)
					alt.RemoveProcessedAlert(alertID)
				} else {
					// unknown error
					appLogger.Error(err)
					alt.RemoveProcessedAlert(alertID)
				}
			}(alertID)
		}
	}()

	appLogger.Info("Starting server")
	fmt.Println(http.ListenAndServe(":8081", mux))

}
