package api

import "net/http"

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Subject  string `json:"subject"`
		Mesasage string `json:"message"`
	}

	var requestPayload mailMessage
	err := app.ReadJson(w, r, &requestPayload)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Data:    requestPayload.Mesasage,
		Subject: requestPayload.Subject,
	}
	err = app.Mailer.SendSMTPMessage(&msg)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Mail sent successfully to " + requestPayload.To,
	}
	app.WriteJson(w, http.StatusAccepted, payload)
}
