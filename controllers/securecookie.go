package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

// pointer to the struct that encodes and decodes
var sCookie *securecookie.SecureCookie

func init() {

	randHashKey := securecookie.GenerateRandomKey(32)
	randBlockKey := securecookie.GenerateRandomKey(32)
	var hashkey = randHashKey
	var blockKey = randBlockKey
	sCookie = securecookie.New(hashkey, blockKey)
	fmt.Println("sCookie has been initialized")
	//env := godotenv.Load("..\\src\\goapp1\\.env")
}

// GetSecureCookie will return global SecureCookie variable
func GetSecureCookie() *securecookie.SecureCookie {
	return sCookie
}

// PrepareCookie will prep a secure cookie with a payload
// params allow you to add a jwt plus its value as well as a csrf and its value
// Options are jwt and _gorilla_csrf
func PrepareCookie(token string, tokenName string) *http.Cookie {

	val := map[string]string{
		tokenName: token,
	}

	var cookie *http.Cookie
	if encodedCookie, err := sCookie.Encode(tokenName, val); err == nil {
		tempCookie := &http.Cookie{
			Name:     tokenName,
			Value:    encodedCookie,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		cookie = tempCookie
	}

	return cookie
}
