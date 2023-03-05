package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/logs/create", app.WriteLog)
	mux.Get("/logs", app.ReadLogs)
	mux.Get("/logs/{id}", app.ReadOne)

	return mux
}
