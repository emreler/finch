package handlers

import (
	"fmt"
	"time"

	"github.com/emreler/finch/channel"
)

// ProcessAlert processes the alert
func (h *Handlers) ProcessAlert(alertID string) error {
	alert, err := h.stg.GetAlert(alertID)

	if err != nil {
		h.logger.Error(err)
		return err
	}

	if alert.Enabled == true && (alert.Schedule.RepeatCount == -1 || alert.Schedule.RepeatCount > 0) {
		h.logger.Info(fmt.Sprintf("Processing %s", alertID))
		if alert.Channel == typeHTTP {
			httpChannel := channel.NewHTTPChannel(h.logger)
			statusCode, err := httpChannel.Notify(alert)

			if err != nil {
				h.logger.Info(fmt.Sprintf("Error while notifying with HTTP channel. %s", err.Error()))
			}

			h.stg.LogProcessAlert(alert, statusCode)

			h.counterChannel <- true
		}

		if alert.Schedule.RepeatCount > 0 {
			alert.Schedule.RepeatCount--

			h.stg.UpdateAlert(alert)
		}
	}

	if alert.Schedule != nil && alert.Schedule.RepeatEvery > 0 {
		h.logger.Info(fmt.Sprintf("Scheduling next alert %d seconds later", alert.Schedule.RepeatEvery))
		h.alt.AddAlert(alertID, time.Duration(alert.Schedule.RepeatEvery)*time.Second)
	}

	return nil
}
