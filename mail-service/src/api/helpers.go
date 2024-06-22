package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Message string      `json:"message"`
	Error   bool        `json:"error"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *Config) ReadJson(
	w http.ResponseWriter,
	r *http.Request,
	data interface{},
) error {
	maxBytes := 1_048_576
	reader := http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}
	return nil
}

func (app *Config) WriteJson(
	w http.ResponseWriter,
	status int,
	data interface{},
	headers ...http.Header,
) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	return err
}

func (app *Config) ErrorJson(
	w http.ResponseWriter,
	err error,
	status ...int,
) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}
	return app.WriteJson(w, statusCode, payload)
}
