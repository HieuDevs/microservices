package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.ReadJson(w, r, &requestPayload)
	if err != nil {
		app.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	// Validate the user agains the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	if user == nil {
		app.ErrorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	//log authentication
	err = app.logRequest("authentication", fmt.Sprintf("User %s logged in", user.Email))
	if err != nil {
		app.ErrorJson(w, err)
		return

	}
	payload := jsonResponse{
		Message: fmt.Sprintf("Welcome %s", user.Email),
		Data:    user,
		Error:   false,
	}
	app.WriteJson(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data
	jsonData, _ := json.MarshalIndent(entry, "", " \t")
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest(
		"POST",
		logServiceURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	return nil
}
