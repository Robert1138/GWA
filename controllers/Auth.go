package controllers

import (
	"fmt"
	u "goapp1/util"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SigningKey Holds the secret for signing jwt
var SigningKey = []byte(os.Getenv("token_secret"))

// JwtMiddleware checks requests that require auth
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuthRequired := []string{"/api/Login", "/api/hello"}
		requestURL := r.URL.Path

		if noAuthRequired[0] == requestURL || noAuthRequired[1] == requestURL {
			next.ServeHTTP(w, r)
			return
		}

		tokenHeader := r.Header.Get("Authorization")
		fmt.Println("JWT from client request", tokenHeader)
		if tokenHeader == "" {
			fmt.Println("no jwt")
			response := u.Message("fail", "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}
		//CreateToken()
		next.ServeHTTP(w, r)
	})
}

// CreateToken call this to create the token using secret and return the string
func CreateToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = true
	claims["name"] = "Bob Smith"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(SigningKey)
	//w.Write([]byte(tokenString))
	fmt.Println(tokenString)
	return tokenString
}
