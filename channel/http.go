package channel

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type HttpChannel struct {
}

func (h *HttpChannel) Notify(url string, data string) error {
	resp, err := http.Post(url, "text/plain", strings.NewReader(data))

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	log.Printf("Response from %s: %s", url, string(body))

	return nil
}
