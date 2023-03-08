package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/logs/create", app.WriteLog)
	mux.Get("/logs", app.ReadLogs)
	mux.Get("/logs/{id}", app.ReadOne)

	return mux
}
