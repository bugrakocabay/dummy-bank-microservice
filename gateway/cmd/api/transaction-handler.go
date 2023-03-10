package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type TransactionRequestPayload struct {
	Action string                   `json:"action"`
	Create CreateTransactionPayload `json:"create,omitempty"`
}

func (app *Config) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	var requestPayload TransactionRequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "create":
		app.createTransactionRequest(w, requestPayload.Create)
	case "get":
		app.getTransactionRequest(w, r)
	case "list":
		app.listTransactionsRequest(w, r)
	default:
		app.errorJSON(w, errors.New(fmt.Sprintf("unknown action type: %s", requestPayload.Action)))
		return
	}
}

type CreateTransactionPayload struct {
	FromAccountID     string         `json:"from_account_id" binding:"required"`
	ToAccountID       string         `json:"to_account_id" binding:"required"`
	TransactionAmount int32          `json:"transaction_amount" binding:"required"`
	Description       sql.NullString `json:"description"`
}

func (app *Config) createTransactionRequest(w http.ResponseWriter, payload CreateTransactionPayload) {
	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPost, "http://account-service/transactions/create", bytes.NewBuffer(jsonData))
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
		app.errorJSON(w, errors.New("error calling transaction service"), response.StatusCode)
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

func (app *Config) getTransactionRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://account-service/transactions/%s", id), strings.NewReader(""))
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
		app.errorJSON(w, errors.New("error calling transaction service"), response.StatusCode)
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

func (app *Config) listTransactionsRequest(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://account-service/transactions"), strings.NewReader(""))

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
		app.errorJSON(w, errors.New("error calling transaction service"), response.StatusCode)
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
