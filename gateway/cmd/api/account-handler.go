package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type AccountRequestPayload struct {
	Action string        `json:"action"`
	Create CreatePayload `json:"create,omitempty"`
	Update UpdatePayload `json:"update,omitempty"`
}

func (app *Config) HandleAccounts(w http.ResponseWriter, r *http.Request) {
	var requestPayload AccountRequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "create":
		app.createAccountRequest(w, r, requestPayload.Create)
	case "get":
		app.getAccountRequest(w, r)
	case "update":
		app.updateAccountRequest(w, requestPayload.Update)
	case "delete":
		app.deleteAccountRequest(w, r)
	default:
		app.errorJSON(w, errors.New(fmt.Sprintf("unknown action type: %s", requestPayload.Action)))
		return
	}
}

// deleteAccountRequest sends an HTTP request to account-service for fetching an existing account
func (app *Config) deleteAccountRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://account-service/accounts/delete/%s", id), strings.NewReader(""))
	if err != nil {
		app.errorJSON(w, err, 500)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, response.StatusCode)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		app.errorJSON(w, errors.New("invalid request"), response.StatusCode)
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling account service"), response.StatusCode)
		return
	}

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		app.errorJSON(w, errors.New("error reading response body"), response.StatusCode)
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "success"
	resp.Data = jsonResponseBody

	app.writeJSON(w, http.StatusOK, resp)
}

// getAccountRequest sends an HTTP request to account-service for fetching an existing account
func (app *Config) getAccountRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://account-service/accounts/%s", id), strings.NewReader(""))
	if err != nil {
		app.errorJSON(w, err, 500)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, response.StatusCode)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		app.errorJSON(w, errors.New("invalid request"), response.StatusCode)
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling account service"), response.StatusCode)
		return
	}

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		app.errorJSON(w, errors.New("error reading response body"), response.StatusCode)
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "success"
	resp.Data = jsonResponseBody

	app.writeJSON(w, http.StatusOK, resp)
}

type UpdatePayload struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int32 `json:"balance" binding:"required"`
}

// updateAccountRequest sends an HTTP request to account-service for updating an existing account
func (app *Config) updateAccountRequest(w http.ResponseWriter, payload UpdatePayload) {
	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPut, "http://account-service/accounts/update", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, 500)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, response.StatusCode)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		app.errorJSON(w, errors.New("invalid request"), response.StatusCode)
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error calling account service"), response.StatusCode)
		return
	}

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		app.errorJSON(w, errors.New("error reading response body"), response.StatusCode)
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "success"
	resp.Data = jsonResponseBody

	app.writeJSON(w, http.StatusOK, resp)
}

type CreatePayload struct {
	Currency string `json:"currency"`
	UserID   any    `json:"user_id"`
}

// createAccountRequest sends an HTTP request to account-service for creating a new account
func (app *Config) createAccountRequest(w http.ResponseWriter, r *http.Request, payload CreatePayload) {
	requestBody := CreatePayload{
		Currency: payload.Currency,
		UserID:   r.Context().Value("user_id"),
	}
	jsonData, _ := json.Marshal(requestBody)

	request, err := http.NewRequest(http.MethodPost, "http://account-service/accounts/create", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, 500)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, response.StatusCode)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		app.errorJSON(w, errors.New("invalid request"), response.StatusCode)
		return
	} else if response.StatusCode != http.StatusCreated {
		app.errorJSON(w, errors.New("error calling account service"), response.StatusCode)
		return
	}

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		app.errorJSON(w, errors.New("error reading response body"), response.StatusCode)
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "success"
	resp.Data = jsonResponseBody

	app.writeJSON(w, http.StatusCreated, resp)
}
