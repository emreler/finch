package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.com/emreler/finch/channel"
)

const (
	TypeHTTP = "http"
)

type Handlers struct {
	stg *Storage
	alt *Alerter
}

func InitHandlers(stg *Storage, alt *Alerter) *Handlers {
	h := &Handlers{stg: stg, alt: alt}

	return h
}

type NewAlertRequest struct {
	Name       string `json:"name"`
	Channel    string `json:"channel"`
	URL        string `json:"url"`
	Data       string `json:"data"`
	AlertDate  string `json:"alertDate"`
	AlertAfter int    `json:"alertAfter"`
}

func (h *Handlers) NewAlert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req := new(NewAlertRequest)

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "invalid format")
			return
		}

		var alertDate time.Time

		if req.AlertAfter != 0 {
			alertDate = time.Now().Add(time.Duration(req.AlertAfter) * time.Second)
		} else if req.AlertDate != "" {
			alertDate, err = time.Parse("2006-01-02T15:04:05Z", req.AlertDate)

			if err != nil {
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return
			}
		}

		alert := &Alert{Name: req.Name, AlertDate: alertDate, Channel: req.Channel, URL: req.URL, Data: req.Data}

		alertID := h.stg.AddAlert(alert)
		seconds := int(alertDate.Sub(time.Now()).Seconds())
		log.Printf("%d seconds later", seconds)

		h.alt.AddAlert(alertID, alertDate)
	}
}

func (h *Handlers) ProcessAlert(alertID string) {
	log.Printf("Getting %s\n", alertID)
	alert := h.stg.GetAlert(alertID)
	log.Printf("Got %+v", alert)

	if alert.Channel == TypeHTTP {
		h := new(channel.HttpChannel)
		h.Notify(alert.URL, alert.Data)
	}
}
