package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, status_msg, data any, others ...string) error {
	message := ""
	if len(others) > 0 {
        message = others[0]
    }
	res := map[string]interface{}{
		"status": status_msg,
        "data":  data,
		"message": message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res)

}


func WriteError(w http.ResponseWriter, status int, err string) {
    WriteJSON(w, status, "error", nil, err)
}


func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return errors.New("request body is missing")
	}
	return json.NewDecoder(r.Body).Decode(payload)

}