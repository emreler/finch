package auth

import (
	"os"
	"testing"
	"time"
)

const secret = "nI1SynC+UUJOi661rkRn614BGCDV2VzzlJKMtFgbpGw="

var tokenString string
var auth *Auth

func TestMain(m *testing.M) {
	auth = NewAuth(secret)
	os.Exit(m.Run())
}

func TestGenerateToken(t *testing.T) {
	exp := time.Now().Add(1 * time.Hour)
	token, err := auth.GenerateToken("123", exp)

	if err != nil {
		t.Error(err)
		return
	}

	tokenString = token

	t.Logf("Created token: %s", tokenString)

}

func TestValidateToken(t *testing.T) {
	err := auth.ValidateToken(tokenString, "123")

	if err != nil {
		t.Error(err)
		return
	}
}
