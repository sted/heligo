package heligo

import (
	"encoding/json"
	"net/http"
)

func writeContentType(w http.ResponseWriter, ct []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = ct
	}
}

func writeJSONContentType(w http.ResponseWriter) {
	writeContentType(w, []string{"application/json", "charset=utf-8"})
}

func writeHTMLContentType(w http.ResponseWriter) {
	writeContentType(w, []string{"text/html", "charset=utf-8"})
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

// WriteJSON write a JSON body with the appropriate header.
func WriteJSON(w http.ResponseWriter, status int, obj any) (int, error) {
	writeJSONContentType(w)
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

func WriteJSONString(w http.ResponseWriter, status int, json []byte) (int, error) {
	writeJSONContentType(w)
	w.WriteHeader(status)
	if !bodyAllowedForStatus(status) {
		return status, nil
	}
	_, err := w.Write(json)
	return status, err
}

func WriteHTMLString(w http.ResponseWriter, status int, html string) (int, error) {
	writeHTMLContentType(w)
	w.WriteHeader(status)
	if !bodyAllowedForStatus(status) {
		return status, nil
	}
	_, err := w.Write([]byte(html))
	return status, err
}

func WriteEmpty(w http.ResponseWriter, status int) (int, error) {
	w.WriteHeader(status)
	return status, nil
}
