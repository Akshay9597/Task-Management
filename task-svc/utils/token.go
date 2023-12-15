package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var Audience = "http://gateway:8080"
var Issuer = "http://users-service:8000"

type (
	AccessToken struct {
		jwt.StandardClaims
		UserId int `json:"user_id"`
	}
	TokenInput struct {
		UserId int
		ExpiresAt int64
	}
)

func generateNewToken(tokenInput TokenInput) AccessToken {
	return AccessToken{
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: tokenInput.ExpiresAt,
			Audience: Audience,
			Issuer: Issuer,
		},
		UserId: tokenInput.UserId,
	}
}

func parseToken(token string) (AccessToken, error) {
	t, _, err := new(jwt.Parser).ParseUnverified(token, &AccessToken{})
	if err != nil {
		return AccessToken{}, err
	}

	claims, ok := t.Claims.(*AccessToken)

	if !ok {
		return AccessToken{}, errors.New("Invalid claims type")
	}

	return *claims, nil
}

