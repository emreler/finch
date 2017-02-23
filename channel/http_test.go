package channel

import (
	"testing"

	"github.com/emreler/finch/models"
)

func TestNotify(t *testing.T) {
	h := &HttpChannel{}

	statusCode, err := h.Notify(&models.Alert{URL: "https://google.com/", Method: methodGet})

	if err != nil {
		t.Error(err)
	}

	t.Logf("status code: %d", statusCode)
}
