package main

import (
	"encoding/json"
	"net/http"
)

// message: Convenience function
func message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// respond: Convenience function
func respond(w http.ResponseWriter, data map[string]interface{}) {

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
