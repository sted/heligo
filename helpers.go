package heligo

import (
	"encoding/json"
	"net/http"
)

// needsClean reports whether the path contains //, /./ or /../ sequences.
func needsClean(p string) bool {
	for i := 0; i < len(p)-1; i++ {
		if p[i] == '/' && (p[i+1] == '/' || p[i+1] == '.') {
			return true
		}
	}
	return false
}

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
	if !bodyAllowedForStatus(status) {
		w.WriteHeader(status)
		return status, nil
	}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return http.StatusInternalServerError, err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, err = w.Write(jsonBytes)
	return status, err
}

// WriteHeader is just for convenience
func WriteHeader(w http.ResponseWriter, status int) (int, error) {
	w.WriteHeader(status)
	return status, nil
}
