package api

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	Message string `json:"message"`
}

type ApiVersion struct {
	Version string `json:"version"`
}

func ApiErrorMsg(w http.ResponseWriter, message string) {
	payload := &ApiError{Message: message}
	data, _ := json.Marshal(payload)
	w.WriteHeader(http.StatusNotAcceptable)
	w.Write(data)
}

func ApiVersionResponse(w http.ResponseWriter, version string) {
	payload := &ApiVersion{Version: version}
	data, _ := json.Marshal(payload)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
