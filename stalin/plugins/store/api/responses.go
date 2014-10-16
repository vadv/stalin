package api

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	Message string `json:"message"`
}

func ApiErrorMsg(w http.ResponseWriter, message string) {
	payload := &ApiError{Message: message}
	data, _ := json.Marshal(payload)
	w.WriteHeader(http.StatusNotAcceptable)
	w.Write(data)
}
