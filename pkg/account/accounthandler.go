package account

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Robert1138/GWA/util"
	lg "github.com/Robert1138/GWA/util/log"

	"github.com/gorilla/mux"
)

var log *lg.StandardLogger

type Passwords struct {
	OldPassword string
	NewPassword string
}

// HTTPRoutes registers endpoints and their appropriate HandlerFuncs to the provided router as well as any subrouters.
func HTTPRoutes(router *mux.Router, newLogger *lg.StandardLogger) {
	log = newLogger
	accountSubrouter := router.PathPrefix("/account").Subrouter()
	accountSubrouter.HandleFunc("/change-password", ChangePassword()).Methods("PUT")
}

// ChangePassword receives the old and new password and changes the user's password to the new one.  this only works for authenticated users
func ChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		var passwordset Passwords
		if err := dec.Decode(&passwordset); err != nil {
			log.Error(errors.New(err.Error() + errors.New(" ChangePassword: accounthandler").Error()))
			util.Respond(w, &util.StandardResponse{Status: 422, Type: "Unprocessable Entity", Error: true, Message: "received incorrect json body"})
		}
		userID, convErr := strconv.Atoi(r.Header.Get("UserID"))
		if convErr != nil {
			log.Error(errors.New(convErr.Error() + errors.New(" ChangePassword: accounthandler").Error()))
			util.Respond(w, &util.StandardResponse{Status: 500, Type: "Internal Server Error", Error: true, Message: convErr.Error()})
		}
		if updateErr := UpdatePassword(userID, passwordset.NewPassword); updateErr != nil {
			log.Error(errors.New(updateErr.Error() + errors.New(" ChangePassword: accounthandler").Error()))
			util.Respond(w, &util.StandardResponse{Status: 404, Type: "Not Found", Error: true, Message: updateErr.Error()})
		}

	}
}
