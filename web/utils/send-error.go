package utils

import "net/http"

func SendError(w http.ResponseWriter, status int, err error /*data interface{}*/) {
	SendJson(w, status, map[string]any{
		"status":  false,
		"message": err.Error(),
		//"data":    data,
	})
}
