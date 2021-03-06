package channel

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/emreler/finch/logger"
	"github.com/emreler/finch/models"
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	contentPlain = "text/plain"
	contentJSON  = "application/json"
	contentForm  = "application/x-www-form-urlencoded"
)

// HTTPChannel implements the http request alert method
type HTTPChannel struct {
	logger logger.InfoErrorLogger
}

// NewHTTPChannel returns new HTTPChannel struct.
func NewHTTPChannel(logger logger.InfoErrorLogger) *HTTPChannel {
	return &HTTPChannel{logger: logger}
}

func (h *HTTPChannel) Notify(alert *models.Alert) (int, error) {
	ValidMethods := map[string]bool{
		methodGet:  true,
		methodPost: true,
	}

	ValidContentTypes := map[string]bool{
		contentPlain: true,
		contentJSON:  true,
		contentForm:  true,
	}

	if alert.Method == "" {
		alert.Method = methodGet
	}

	if !ValidMethods[alert.Method] {
		return 0, fmt.Errorf("Invalid method %s", alert.Method)
	}

	if alert.ContentType == "" {
		alert.ContentType = contentPlain
	}

	if !ValidContentTypes[alert.ContentType] {
		return 0, fmt.Errorf("Invalid contentType %s", alert.ContentType)
	}

	var resp *http.Response
	var err error

	httpClient := http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
		Timeout:   time.Second * 10,
	}

	if alert.Method == methodGet {
		resp, err = httpClient.Get(alert.URL)
	} else if alert.Method == methodPost {
		resp, err = httpClient.Post(alert.URL, alert.ContentType, strings.NewReader(alert.Data))
	}

	if err != nil {
		return 0, err
	}

	resp.Body.Close()

	h.logger.Info(fmt.Sprintf("Response for %s request to %s: %d", alert.Method, alert.URL, resp.StatusCode))

	return resp.StatusCode, nil
}
