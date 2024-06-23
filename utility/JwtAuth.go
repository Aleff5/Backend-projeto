package utility

import (
	"context"
	"errors"
	"net/http"
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		claims := &jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET_KEY), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", (*claims)["username"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
