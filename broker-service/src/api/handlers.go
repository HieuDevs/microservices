package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayLoad struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit Broker API",
	}
	_ = app.WriteJson(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayLoad
	err := app.ReadJson(w, r, &requestPayload)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.log(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.ErrorJson(w, errors.New("invalid action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json we'll send to authentication service
	jsonData, _ := json.MarshalIndent(a, "", " \t")
	// call the service
	request, err := http.NewRequest(
		"POST",
		"http://authentication-service/authenticate",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	defer response.Body.Close()
	// make sture we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJson(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.ErrorJson(w, errors.New("authentication service error"))
		return
	}
	// read the response
	var jsonFromAuthService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromAuthService)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	if jsonFromAuthService.Error {
		app.ErrorJson(w, err, http.StatusUnauthorized)
		return
	}
	// send the response back to the client
	var payloadResponse jsonResponse
	payloadResponse.Message = "Authenticated"
	payloadResponse.Data = jsonFromAuthService.Data
	payloadResponse.Error = false
	_ = app.WriteJson(w, http.StatusAccepted, payloadResponse)
}

func (app *Config) log(w http.ResponseWriter, l LogPayload) {
	// create json we'll send to logging service
	jsonData, _ := json.MarshalIndent(l, "", " \t")
	// call the service
	request, err := http.NewRequest(
		"POST",
		"http://logger-service/log",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	defer response.Body.Close()
	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		app.ErrorJson(w, errors.New("logging service error"))
		return
	}
	// read the response
	var jsonFromLoggingService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromLoggingService)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	if jsonFromLoggingService.Error {
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	// send the response back to the client
	var payloadResponse jsonResponse
	payloadResponse.Message = "Logged"
	payloadResponse.Data = jsonFromLoggingService.Data
	payloadResponse.Error = false
	_ = app.WriteJson(w, http.StatusAccepted, payloadResponse)
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	// create json we'll send to logging service
	jsonData, _ := json.MarshalIndent(m, "", " \t")
	// call the service
	request, err := http.NewRequest(
		"POST",
		"http://mail-service/send",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	defer response.Body.Close()
	// make sure we get back the correct status code
	if response.StatusCode != http.StatusAccepted {
		app.ErrorJson(w, errors.New("mail service error"))
		return
	}
	// read the response
	var jsonFromMailService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromMailService)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	if jsonFromMailService.Error {
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	// send the response back to the client
	var payloadResponse jsonResponse
	payloadResponse.Message = "Mail sent successfully to " + m.To
	payloadResponse.Data = jsonFromMailService.Data
	payloadResponse.Error = false
	_ = app.WriteJson(w, http.StatusAccepted, payloadResponse)
}
