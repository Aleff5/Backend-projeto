package utility

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JWT_SECRET_KEY = "mykey"

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 5)

	claims := jwt.MapClaims{
		"username": username,
		"expires":  expirationTime.Unix(),
	}
	//gera token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//assina token
	return token.SignedString([]byte(JWT_SECRET_KEY))
}
