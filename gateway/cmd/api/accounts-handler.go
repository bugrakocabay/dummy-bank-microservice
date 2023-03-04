package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string        `json:"action"`
	Create CreatePayload `json:"create,omitempty"`
}

func (app *Config) HandleAccounts(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "create":
		app.createAccountRequest(w, requestPayload.Create)
	}
}

type CreatePayload struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// createAccountRequest sends an HTTP request to account-service for creating a new account
func (app *Config) createAccountRequest(w http.ResponseWriter, payload CreatePayload) {
	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPost, "http://account-service/accounts/create", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		app.errorJSON(w, errors.New("error reading response body"))
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "success"
	resp.Data = jsonResponseBody

	app.writeJSON(w, http.StatusCreated, resp)
}
