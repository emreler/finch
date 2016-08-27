package channel

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type HttpChannel struct {
}

func (h *HttpChannel) Notify(url string, data string) {
	resp, _ := http.Post(url, "text/plain", strings.NewReader(data))
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	log.Printf("Response from %s: %s", url, string(body))
}
