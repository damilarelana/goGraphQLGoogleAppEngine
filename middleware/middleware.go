package middleware

import (
	"encoding/json"
	"net/http"
)

// ResponseError endpoint handler
func ResponseError(w http.ResponseWriter, errMsg string, errCode int) {
	w.Header().Set("Content-Type", "application/json: charset=UFT-8") // set the content header type
	w.WriteHeader(errCode)
	json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
}

// ResponseJSON endpoint handler
func ResponseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json: charset=UFT-8") // set the content header type
	json.NewEncoder(w).Encode(data)
}
