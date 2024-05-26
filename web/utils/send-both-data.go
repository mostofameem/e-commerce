package utils

import (
	"net/http"
)

func SendBothData(w http.ResponseWriter, pageinfo interface{}, data interface{}) {
	SendJson(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": pageinfo,
		"data":    data,
	})
}
