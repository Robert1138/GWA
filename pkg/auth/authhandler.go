package auth

import (
	"encoding/json"
	"fmt"
	u "goapp1/pkg/user"
	"goapp1/util"
	"goapp1/util/jwt"
	lg "goapp1/util/log"
	"net/http"

	"github.com/gorilla/mux"
)

var log *lg.StandardLogger

type userLogin struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

// HTTPRoutes registers endpoints and their appropriate HandlerFuncs to the provided router as well as any subrouters.
func HTTPRoutes(router *mux.Router, newLogger *lg.StandardLogger) {
	log = newLogger
	router.HandleFunc("/login", Login()).Methods("POST")
	router.HandleFunc("/register", Register()).Methods("POST")
}

// Login checks provided credentials, if valid a token is generated and a secure cookie is sent, else send back reasons why not
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user userLogin
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// when user credentials are identified we can set token and send cookie with jwt
		// boolean, User returned
		fmt.Println("User sent password: " + user.UserName + " " + user.Password)
		userCred, isValid := ValidUserPassword(user.UserName, user.Password)
		if isValid {
			log.Info("User valid login " + user.UserName)
			fmt.Println(userCred.UserID)
			http.SetCookie(w, jwt.PrepareCookie(jwt.CreateToken(w, userCred), "jwt"))
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity) // send a 422 -- user they sent does not exist
			util.Response(w, util.Message("422", "Unprocessable Entity - check what you're sending - user might not exist"))
			//w.Write([]byte("422 - Unprocessable Entity - check what you're sending - user might not exist"))
		}
	}
}

// Register will
func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser u.User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		//newUser.Password, err = hashPassword(newUser.Password)
		if err != nil {
			fmt.Println("password couldnt be encrypted")
		}
		userAdded, err := AddUser(newUser)

		if err != nil {
			fmt.Println(err)
			if err == util.ErrUserExists {
				fmt.Println("user exists")
				w.WriteHeader(http.StatusConflict)
				util.Response(w, util.Message("409", "user with this UserName already exists"))
				w.Write([]byte("409 - user with this UserName already exists"))
			} else if err == util.ErrPasswordInvalidFormat {
				w.WriteHeader(http.StatusUnauthorized)
				util.Response(w, util.Message("401", "password too short"))
				w.Write([]byte("401 - password too short"))
			}
		}
		http.SetCookie(w, jwt.PrepareCookie(jwt.CreateToken(w, newUser), "jwt"))
		fmt.Println(userAdded)
	}
}
