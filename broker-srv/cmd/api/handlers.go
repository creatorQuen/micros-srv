package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (a *Config) Broker(w http.ResponseWriter, r *http.Request) {
	body := jsonResponse{
		Error:   false,
		Message: "Puchnch broker",
	}

	out, _ := json.MarshalIndent(body, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}
