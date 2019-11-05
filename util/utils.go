package util

import (
	"encoding/json"
	"net/http"
)

// sends back a string - SUCCESS or FAILURE and another string indicating why
func Message(status string, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// sends a reponse with the data payload
func Response(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
