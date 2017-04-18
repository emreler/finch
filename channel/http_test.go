package channel

import (
	"testing"

	"github.com/emreler/finch/logger"
	"github.com/emreler/finch/models"
)

func TestNotify(t *testing.T) {
	mockLogger := &logger.MockLogger{}
	h := NewHTTPChannel(mockLogger)

	statusCode, err := h.Notify(&models.Alert{URL: "http://example.com/", Method: methodGet})

	if err != nil {
		t.Error(err)
	}

	t.Logf("status code: %d", statusCode)
}
