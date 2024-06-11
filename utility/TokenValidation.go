package utility

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func Validation(tokenstr string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenstr, func(t *jwt.Token) (interface{}, error) {
		//check the jwt method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWT_SECRET_KEY), nil
	})
	if err != nil {
		return nil, errors.New("Invalid Token")
	}
	//check if the token is valid
	if !token.Valid {
		return nil, errors.New("Invalid Token")
	}
	//get the claims of the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid Token")

	}
	return claims, nil
}
