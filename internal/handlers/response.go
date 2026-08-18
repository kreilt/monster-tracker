package handlers

import (
	"encoding/json"
	"net/http"
)

func writeError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
