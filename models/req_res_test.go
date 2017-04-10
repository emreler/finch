package models

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCreateUserRequest(t *testing.T) {
	jsonBlob := []byte(`{
			"name": "foo",
			"email": "bartar.com"
	}`)

	req := &CreateUserRequest{}
	err := json.Unmarshal(jsonBlob, req)

	if err != nil {
		t.Error(err)
	}

	err = req.Validate()

	if err == nil {
		t.Error(fmt.Errorf("it should return error for invalid email address"))
	}
}

func TestCreateAlertRequest(t *testing.T) {
	jsonBlob := []byte(`{
		"alertAfter": 2,
		"channel": "asd",
		"contentType": "text/plain",
		"method": "GET",
		"repeatEvery": 1,
		"url": "http://example.com"
	}`)

	req := &CreateAlertRequest{}
	err := json.Unmarshal(jsonBlob, req)

	if err != nil {
		t.Error(err)
	}

	err = req.Validate()

	if err == nil {
		t.Error(fmt.Errorf("it should return error for invalid channel name"))
	}
}
