package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	accountServiceURL = "http://account-service"
	maxBytes          = 10485376
)

type AccountRequestPayload struct {
	Action  string        `json:"action"`
	Create  CreatePayload `json:"create,omitempty"`
	Update  UpdatePayload `json:"update,omitempty"`
	Balance AddBalance    `json:"balance,omitempty"`
}

type accountResponse struct {
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	AccountID string    `json:"account_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (app *Config) HandleAccounts(w http.ResponseWriter, r *http.Request) {
	var requestPayload AccountRequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, "HandleAccounts", err)
	}

	switch requestPayload.Action {
	case "create":
		app.createAccountRequest(w, r, requestPayload.Create)
	case "get":
		app.getAccountRequest(w, r)
	case "update":
		app.updateAccountRequest(w, r, requestPayload.Update)
	case "delete":
		app.deleteAccountRequest(w, r)
	case "list":
		app.listAccountRequest(w, r)
	case "balance":
		app.addBalanceRequest(w, r, requestPayload.Balance)
	default:
		app.errorJSON(w, "HandleAccounts", errors.New(fmt.Sprintf("unknown action type: %s", requestPayload.Action)))
	}
}

func (app *Config) deleteAccountRequest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "account_id")
	userID := fmt.Sprintf("%v", r.Context().Value("user_id"))

	// Check if the user is authorized to delete the account
	accountUserID, err := getAccountUserID(id)
	if err != nil {
		return app.errorJSON(w, "deleteAccountRequest", err, http.StatusInternalServerError)
	}
	if accountUserID != userID {
		return app.errorJSON(w, "deleteAccountRequest", errors.New("you are not authorized to delete this account"), http.StatusForbidden)
	}

	// Send DELETE request to account service
	reqURL := fmt.Sprintf("%s/accounts/delete/%s", accountServiceURL, id)
	req, err := http.NewRequest(http.MethodDelete, reqURL, nil)
	if err != nil {
		return app.errorJSON(w, "deleteAccountRequest", err, http.StatusInternalServerError)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return app.errorJSON(w, "deleteAccountRequest", err, resp.StatusCode)
	}
	defer resp.Body.Close()

	// Decode JSON response
	var respBody struct {
		Error   bool        `json:"error"`
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	}
	err = json.NewDecoder(http.MaxBytesReader(w, resp.Body, maxBytes)).Decode(&respBody)
	if err != nil {
		return app.errorJSON(w, "deleteAccountRequest", err, resp.StatusCode)
	}

	// Write JSON response
	statusCode := resp.StatusCode
	if statusCode != http.StatusOK {
		respBody.Message = "fail"
	}
	return app.writeJSON(w, "deleteAccountRequest", statusCode, respBody)
}

// getAccountRequest sends an HTTP request to account-service for fetching an existing account
func (app *Config) getAccountRequest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "account_id")

	reqURL := fmt.Sprintf("%s/accounts/%s", accountServiceURL, id)
	request, err := http.NewRequest(http.MethodGet, reqURL, strings.NewReader(""))
	if err != nil {
		return app.errorJSON(w, "getAccountRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "getAccountRequest", err, response.StatusCode)

	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody accountResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "getAccountRequest", errors.New("error reading response body"), response.StatusCode)
	}

	idInHeader := r.Context().Value("user_id")
	if jsonResponseBody.UserID != idInHeader {
		return app.errorJSON(w, "getAccountRequest", errors.New("this is not yours"), 403)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "getAccountRequest", response.StatusCode, resp)
}

type UpdatePayload struct {
	AccountID string `json:"account_id" binding:"required"`
	Balance   int32  `json:"balance" binding:"required"`
}

// updateAccountRequest sends an HTTP request to account-service for updating an existing account
func (app *Config) updateAccountRequest(w http.ResponseWriter, r *http.Request, payload UpdatePayload) error {
	jsonData, _ := json.Marshal(payload)

	idInHeader := fmt.Sprintf("%v", r.Context().Value("user_id"))
	userID, err := getAccountUserID(payload.AccountID)
	if userID != idInHeader {
		return app.errorJSON(w, "updateAccountRequest", errors.New("this is not yours"), 403)
	}

	reqURL := fmt.Sprintf("%s/accounts/update", accountServiceURL)
	request, err := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return app.errorJSON(w, "updateAccountRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.sendErrorLog("updateAccountRequest", Log{
			StatusCode: response.StatusCode,
			Message:    err.Error(),
		})
		return app.errorJSON(w, "updateAccountRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody accountResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "updateAccountRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "updateAccountRequest", response.StatusCode, resp)
}

type CreatePayload struct {
	Currency string `json:"currency"`
	UserID   any    `json:"user_id"`
}

// createAccountRequest sends an HTTP request to account-service for creating a new account
func (app *Config) createAccountRequest(w http.ResponseWriter, r *http.Request, payload CreatePayload) error {
	requestBody := CreatePayload{
		Currency: payload.Currency,
		UserID:   r.Context().Value("user_id"),
	}
	jsonData, _ := json.Marshal(requestBody)

	reqURL := fmt.Sprintf("%s/accounts/create", accountServiceURL)
	request, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return app.errorJSON(w, "createAccountRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "createAccountRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "createAccountRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "createAccountRequest", response.StatusCode, resp)
}

func (app *Config) listAccountRequest(w http.ResponseWriter, r *http.Request) error {
	reqURL := fmt.Sprintf("%s/accounts", accountServiceURL)
	request, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return app.errorJSON(w, "listAccountRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "listAccountRequest", err, response.StatusCode)

	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "listAccountRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "listAccountRequest", response.StatusCode, resp)
}

type AddBalance struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

func (app *Config) addBalanceRequest(w http.ResponseWriter, r *http.Request, payload AddBalance) error {
	requestBody := AddBalance{
		AccountID: payload.AccountID,
		Amount:    payload.Amount,
	}
	jsonData, _ := json.Marshal(requestBody)

	reqURL := fmt.Sprintf("%s/accounts/add-balance", accountServiceURL)
	request, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return app.errorJSON(w, "addBalanceRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "addBalanceRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody accountResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "createAccountRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "addBalanceRequest", response.StatusCode, resp)
}
