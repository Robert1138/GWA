package controllers

import (
	"encoding/json"
	"fmt"
	u "goapp1/util"
	"net/http"
)

type userLogin struct {
	Name     string
	Password string
}

// GetMessage returns a simple response
func GetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u.Response(w, u.Message("success", "it worked"))
	}
}

// Login will take credentials, verify and send a jwt in the response, TODO use secure cookies
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO check credentials THEN return response
		// userLoginEvent := userLoginCredentials{}
		//rBody, err := ioutil.ReadAll(r.Body)
		//fmt.Println(rBody)
		// Will hold user info from login post
		var user userLogin

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Do something with the Person struct...
		fmt.Println("User: " + user.Name + " " + user.Password)

		//json.Unmarshal(rBody, &userLoginEvent)
		//var userInfo []string

		//u.Response(w, map[string]interface{}{"jwt": CreateToken(arr)})
		//fmt.Println(w.Header)
		w.Header().Set("X-CSRF-Token", "bleh")
		http.SetCookie(w, PrepareCookie(CreateToken(user), "jwt"))

	}
}
