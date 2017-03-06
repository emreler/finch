package auth

import (
	"os"
	"testing"
	"time"
)

const secret = "nI1SynC+UUJOi661rkRn614BGCDV2VzzlJKMtFgbpGw="
const userID = "12345"

var tokenString string
var auth *Auth

func TestMain(m *testing.M) {
	auth = NewAuth(secret)
	os.Exit(m.Run())
}

func TestGenerateToken(t *testing.T) {
	exp := time.Now().Add(1 * time.Hour)
	token, err := auth.GenerateToken(userID, exp)

	if err != nil {
		t.Error(err)
		return
	}

	tokenString = token

	t.Logf("Created token: %s", tokenString)

}

func TestValidateToken(t *testing.T) {
	claimedUserID, err := auth.ValidateToken(tokenString)

	if err != nil {
		t.Error(err)
		return
	}

	if claimedUserID != userID {
		t.Errorf("UserID %s in token is not matching with %s", claimedUserID, userID)
	}
}
