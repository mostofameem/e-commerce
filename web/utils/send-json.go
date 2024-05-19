package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func SendJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")

	str, err := json.Marshal(data)

	if err != nil {
		err = errors.New("failed converting into json")
		SendError(w, status, err, data)
		return
	}
	w.WriteHeader(status)
	w.Write(str)
}
