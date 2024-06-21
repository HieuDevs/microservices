package api

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit Broker API",
	}
	_ = app.WriteJson(w, http.StatusOK, payload)
}
