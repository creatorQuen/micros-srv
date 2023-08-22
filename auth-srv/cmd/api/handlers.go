package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(payload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid payload"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatch(payload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid payload"), http.StatusBadRequest)
	}

	answerPayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, answerPayload)
}
