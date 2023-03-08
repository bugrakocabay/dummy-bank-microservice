package main

import (
	"net/http"

	"github.com/bugrakocabay/dummy-bank-microservice/logger-service/data"
	"github.com/go-chi/chi/v5"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "success",
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *Config) ReadLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := app.Models.LogEntry.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    logs,
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Config) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	record, err := app.Models.LogEntry.GetOne(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    record,
	}

	app.writeJSON(w, http.StatusOK, resp)
}
