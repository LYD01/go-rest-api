// JSON encode/decode + error helpers

package respond

import (
	"encoding/json"
	"net/http"
)


func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if (data != nil) {
		json.NewEncoder(w).Encode(data)
	}
}

type errorBody struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, errorBody{Error: msg})
}

func Decode(r *http.Request, dest any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dest)
}
