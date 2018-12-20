package main

import (
	"encoding/json"
	"net/http"
)

//Message: Convenience function
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond: Convenience function
func Respond(w http.ResponseWriter, data map[string]interface{}) {

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
