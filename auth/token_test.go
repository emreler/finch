package auth

import (
	"testing"
	"time"
)

const secret = "nI1SynC+UUJOi661rkRn614BGCDV2VzzlJKMtFgbpGw="

var tokenString string

func TestGenerateToken(t *testing.T) {
	exp := time.Now().Add(1 * time.Hour)
	token, err := GenerateToken("123", exp, secret)

	if err != nil {
		t.Error(err)
		return
	}

	tokenString = token

	t.Logf("Created token: %s", tokenString)

}

func TestValidateToken(t *testing.T) {
	err := ValidateToken(tokenString, "123", secret)

	if err != nil {
		t.Error(err)
		return
	}
}
