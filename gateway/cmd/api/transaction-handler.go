package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type TransactionRequestPayload struct {
	Action string                   `json:"action"`
	Create CreateTransactionPayload `json:"create,omitempty"`
}

func (app *Config) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	var requestPayload TransactionRequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		if err = app.errorJSON(w, "HandleTransactions", err); err != nil {
			return
		}
		return
	}

	switch requestPayload.Action {
	case "create":
		app.createTransactionRequest(w, r, requestPayload.Create)
	case "get":
		app.getTransactionRequest(w, r)
	case "list":
		app.listTransactionsRequest(w, r)
	default:
		if err = app.errorJSON(w, "HandleTransactions", errors.New(fmt.Sprintf("unknown action type: %s", requestPayload.Action))); err != nil {
			return
		}
		return
	}
}

type CreateTransactionPayload struct {
	FromAccountID     string         `json:"from_account_id" binding:"required"`
	ToAccountID       string         `json:"to_account_id" binding:"required"`
	TransactionAmount float64        `json:"transaction_amount" binding:"required"`
	Description       sql.NullString `json:"description"`
}

func (app *Config) createTransactionRequest(w http.ResponseWriter, r *http.Request, payload CreateTransactionPayload) error {
	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPost, "http://account-service/transactions/create", bytes.NewBuffer(jsonData))
	if err != nil {
		return app.errorJSON(w, "createTransactionRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "createTransactionRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "createTransactionRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "createTransactionRequest", response.StatusCode, resp)
}

func (app *Config) getTransactionRequest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "transaction_id")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://account-service/transactions/%s", id), nil)
	if err != nil {
		return app.errorJSON(w, "getTransactionRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "getTransactionRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "getTransactionRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "getTransactionRequest", response.StatusCode, resp)
}

func (app *Config) listTransactionsRequest(w http.ResponseWriter, r *http.Request) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://account-service/transactions"), strings.NewReader(""))

	if err != nil {
		return app.errorJSON(w, "listTransactionRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "listTransactionRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "listTransactionRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "listTransactionsRequest", response.StatusCode, resp)
}
