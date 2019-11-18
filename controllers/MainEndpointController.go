package controllers

import (
	u "goapp1/util"
	"net/http"
)

// GetMessage returns a simple response
func GetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u.Response(w, u.Message("success", "it worked"))
	}
}

// Login will take credentials, verify and send a jwt in the response
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u.Response(w, map[string]interface{}{"jwt": CreateToken()})
	}
}

// GetFavoriteColor returns favorite colors
func GetFavoriteColor(user string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
