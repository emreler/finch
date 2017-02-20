package channel

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/emreler/finch/models"
)

const (
	methodGet    = "GET"
	methodPost   = "POST"
	contentPlain = "text/plain"
	contentJSON  = "application/json"
	contentForm  = "application/x-www-form-urlencoded"
)

type HttpChannel struct {
}

func (h *HttpChannel) Notify(alert *models.Alert) error {
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
		return fmt.Errorf("Invalid method %s", alert.Method)
	}

	if alert.ContentType == "" {
		alert.ContentType = contentPlain
	}

	if !ValidContentTypes[alert.ContentType] {
		return fmt.Errorf("Invalid contentType %s", alert.ContentType)
	}

	var resp *http.Response
	var err error

	if alert.Method == methodGet {
		resp, err = http.Get(alert.URL)
	} else if alert.Method == methodPost {
		resp, err = http.Post(alert.URL, alert.ContentType, strings.NewReader(alert.Data))
	}

	if err != nil {
		return err
	}

	// defer resp.Body.Close()
	//
	// var body []byte
	// if body, err = ioutil.ReadAll(resp.Body); err != nil {
	// 	return err
	// }

	log.Printf("Response for %s request to %s: %d", alert.Method, alert.URL, resp.StatusCode)

	return nil
}
