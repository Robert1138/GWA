package csrftest

import (
	"encoding/json"
	"goapp1/util"
	lg "goapp1/util/log"
	"net/http"

	"github.com/gorilla/mux"
)

var log *lg.StandardLogger

// HTTPRoutes registers endpoints and their appropriate HandlerFuncs to the provided router as well as any subrouters.
func HTTPRoutes(router *mux.Router, newLogger *lg.StandardLogger) {
	log = newLogger
	router.HandleFunc("/stuff", GetMessage()).Methods("GET")   // GET placeholder to get csrf token
	router.HandleFunc("/CSRF", CsrfPostTest()).Methods("POST") // POST placeholder that requires a csrf token in the header
}

// CsrfPostTest tests post with csrf protection
func CsrfPostTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("POST WENT THROUGH WITH VALID CSRF")

	}
}

// GetMessage returns a simple response -- placeholder to get csrf token
func GetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.Response(w, util.Message("success", "it worked"))
	}
}
