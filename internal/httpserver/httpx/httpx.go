package httpx

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return e.Message
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, map[string]string{"error": message})
}

func RequireMethod(r *http.Request, method string) error {
	if r.Method != method {
		return &Error{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "method not allowed",
		}
	}
	return nil
}

func RequiredQuery(r *http.Request, key string) (string, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return "", &Error{
			StatusCode: http.StatusBadRequest,
			Message:    "missing " + key + " query parameter",
		}
	}
	return value, nil
}
