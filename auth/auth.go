package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	secret string
}

func NewAuth(secret string) *Auth {
	return &Auth{secret}
}

func (a *Auth) GenerateToken(userID string, expire time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    expire.Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Auth) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("Invalid token signature")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("Invalid token body")
	}

	exp := time.Unix(int64(claims["exp"].(float64)), 0)

	if !exp.After(time.Now()) {
		return "", fmt.Errorf("Invalid or expired token")
	}

	return claims["userID"].(string), nil
}
