package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *Config) Router() http.Handler {
	mux := chi.NewRouter()
	mux.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"http://*", "https://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: true,
				MaxAge:           300,
			},
		),
	)
	mux.Post("/authenticate", app.Authenticate)
	return mux
}
