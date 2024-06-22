package api

import (
	"logger-service/src/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	//read json into variable
	var payload JSONPayload
	_ = app.ReadJson(w, r, &payload)

	//write to log
	event := data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}
	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	//return success
	resp := jsonResponse{
		Error:   false,
		Message: "Log created",
	}
	app.WriteJson(w, http.StatusAccepted, resp)
}
