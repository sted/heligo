package heligo

import (
	"encoding/json"
	"net/http"
)

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

// WriteJSON writes a JSON body with the appropriate header.
func WriteJSON(w http.ResponseWriter, status int, obj any) (int, error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if !bodyAllowedForStatus(status) {
		return status, nil
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return status, err
	}
	_, err = w.Write(jsonBytes)
	return status, err
}

// WriteHeader is just for convenience
func WriteHeader(w http.ResponseWriter, status int) (int, error) {
	w.WriteHeader(status)
	return status, nil
}
