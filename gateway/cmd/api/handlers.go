package main

import (
	"log"
	"net/http"
)

func (app *Config) Gateway(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "You hit the gateway!",
	}

	if err := app.writeJSON(w, http.StatusOK, payload); err != nil {
		log.Printf("error while sending json response: %v", err)
	}
}
