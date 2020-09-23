package account

import (
	"encoding/json"
	"errors"
	"goapp1/util"
	lg "goapp1/util/log"
	"net/http"

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
		err := dec.Decode(&passwordset)
		if err != nil {
			log.Error(errors.New(err.Error() + errors.New(" ChangePassword: accounthandler").Error()))
			util.Respond(w, &util.StandardResponse{Status: 422, Type: "Unprocessable Entity", Error: true, Message: "received incorrect json body"})
		}
		//updated := UpdatePassword()

	}
}
