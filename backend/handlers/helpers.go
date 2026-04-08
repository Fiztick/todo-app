package handlers

import (
	"encoding/json"
	"net/http"
)

func respondJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJson(w, status, map[string]string{"error": message})
}
