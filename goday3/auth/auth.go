package auth

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	t, err := token()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": t,
	})
}

var Token = token

func token() (string, error) {
	mySigningKey := []byte("drowssap")

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}
