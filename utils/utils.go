package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)

}


func WriteError(w http.ResponseWriter, status int, err error) {
	errorData := map[string]string{
		"status": "error",
		"message": err.Error(),
	}
    WriteJSON(w, status, errorData)
}