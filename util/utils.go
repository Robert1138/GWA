package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// StandardResponse is a struct that holds information about the response to a http request
// Status: response code, Type: description of response code, Error: bool if error, Message: Info pertaining to response
type StandardResponse struct {
	Status  int
	Type    string
	Error   bool
	Message string
}

// Message sends back a string - SUCCESS or FAILURE and another string indicating why
func Message(status string, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Response sends a reponse with the data payload
func Response(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Respond(w http.ResponseWriter, data *StandardResponse) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.Status)
	json.NewEncoder(w).Encode(data)
}

func HashPassword(password string) (string, error) {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(passBytes), err
}

func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
