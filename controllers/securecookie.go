package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
)

// pointer to the struct that encodes and decodes
var sCookie *securecookie.SecureCookie

func init() {

	var hashKey = []byte(os.Getenv("hashkey"))
	var blockKey = []byte(os.Getenv("blockkey"))
	sCookie = securecookie.New(hashKey, blockKey)
	fmt.Println("sCookie has been initialized")
	fmt.Println(os.Getenv("hashkey"))

	// fmt.Println(base64.StdEncoding.EncodeToString([]byte(securecookie.GenerateRandomKey(32))))
	// fmt.Println(base64.StdEncoding.EncodeToString([]byte(securecookie.GenerateRandomKey(32))))
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
			HttpOnly: true,
		}
		cookie = tempCookie
	}

	return cookie
}
