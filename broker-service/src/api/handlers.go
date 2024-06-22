package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"my-broker/src/event"
	"my-broker/src/logs"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		// app.log(w, requestPayload.Log)
		// app.logEventViaRabbit(w, requestPayload.Log)
		app.logEventViaGRPC(w, requestPayload.Log)
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

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	var payloadResponse jsonResponse
	payloadResponse.Message = "Logged via RabbitMQ"
	payloadResponse.Data = l.Data
	payloadResponse.Error = false
	_ = app.WriteJson(w, http.StatusAccepted, payloadResponse)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	j, _ := json.MarshalIndent(LogPayload{
		Name: name,
		Data: msg,
	}, "", " \t")
	return emitter.Push(string(j), "log.INFO")
}

func (app *Config) logEventViaGRPC(w http.ResponseWriter, l LogPayload) {

	conn, err := grpc.NewClient(
		"logger-service:50001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	defer conn.Close()
	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: l.Name,
			Data: l.Data,
		},
	})
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	var payloadResponse jsonResponse
	payloadResponse.Message = "Logged via gRPC"
	payloadResponse.Data = l.Data
	payloadResponse.Error = false
	_ = app.WriteJson(w, http.StatusAccepted, payloadResponse)
}
