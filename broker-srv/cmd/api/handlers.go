package main

import (
	"net/http"
)

func (c *Config) Broker(w http.ResponseWriter, r *http.Request) {
	body := jsonResponse{
		Error:   false,
		Message: "Throght broker",
	}

	_ = c.writeJSON(w, http.StatusOK, body)
}
