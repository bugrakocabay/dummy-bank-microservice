package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/logs/create-error", app.WriteErrorLog)
	mux.Post("/logs/create-request", app.WriteRequestLog)
	mux.Get("/logs", app.ReadLogs)
	mux.Get("/logs/{id}", app.ReadOne)

	return mux
}
