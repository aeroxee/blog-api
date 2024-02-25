package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserAuth struct{}

type Credential struct {
	UserID int `json:"user_id"`
}

type Claims struct {
	Credential
	jwt.RegisteredClaims
}

var secretKey = []byte("secret key")

func GetToken(credential Credential) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := Claims{
		Credential: credential,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(secretKey)
	return token, err
}

func VerifyToken(token string) (Claims, error) {
	var claims Claims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return claims, err
	}

	if !jwtToken.Valid {
		return claims, errors.New("token anda sudah kadaluarsa")
	}
	return claims, nil
}
