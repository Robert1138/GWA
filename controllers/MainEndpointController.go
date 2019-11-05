package controllers

import (
	u "goapp1/util"
	"net/http"
)

func GetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u.Response(w, u.Message("success", "it worked"))
	}
}
