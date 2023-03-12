package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type UserRequestPayload struct {
	Action string            `json:"action"`
	Create CreateUserPayload `json:"create,omitempty"`
	Login  LoginUserPayload  `json:"login,omitempty"`
}

func (app *Config) HandleUsers(w http.ResponseWriter, r *http.Request) {
	var requestPayload UserRequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		if err = app.errorJSON(w, err); err != nil {
			return
		}
		return
	}

	switch requestPayload.Action {
	case "create":
		app.createUserRequest(w, requestPayload.Create)
	case "get":
		app.getUserRequest(w, r)
	case "login":
		app.loginUserRequest(w, requestPayload.Login)
	default:
		if err = app.errorJSON(w, errors.New(fmt.Sprintf("unknown action type: %s", requestPayload.Action))); err != nil {
			return
		}
		return
	}
}

type LoginUserPayload struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (app *Config) loginUserRequest(w http.ResponseWriter, payload LoginUserPayload) {
	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPost, "http://user-service/users/login", bytes.NewBuffer(jsonData))
	if err != nil {
		if err = app.errorJSON(w, err, 500); err != nil {
			return
		}
		app.sendErrorLog("loginUserRequest", errorLog{
			StatusCode: 500,
			Message:    err,
		})
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		if err = app.errorJSON(w, err, response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("loginUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    err,
		})
		return
	}
	defer response.Body.Close()

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)

	if response.StatusCode == http.StatusBadRequest {
		if err = app.errorJSON(w, errors.New("invalid request"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("loginUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    jsonResponseBody,
		})
		return
	} else if response.StatusCode != http.StatusOK {
		if err = app.errorJSON(w, errors.New("error calling user service"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("loginUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    jsonResponseBody,
		})
		return
	}

	if err != nil {
		if err = app.errorJSON(w, errors.New("error reading response body"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("loginUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    err,
		})
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "success"
	resp.Data = jsonResponseBody

	if err = app.writeJSON(w, http.StatusOK, resp); err != nil {
		return
	}
}

type CreateUserPayload struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (app *Config) createUserRequest(w http.ResponseWriter, payload CreateUserPayload) {
	jsonData, _ := json.Marshal(payload)

	request, err := http.NewRequest(http.MethodPost, "http://user-service/users/create", bytes.NewBuffer(jsonData))
	if err != nil {
		if err = app.errorJSON(w, err, 500); err != nil {
			return
		}
		app.sendErrorLog("createUserRequest", errorLog{
			StatusCode: 500,
			Message:    err,
		})
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		if err = app.errorJSON(w, err, response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("createUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    err,
		})
		return
	}
	defer response.Body.Close()

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)

	if response.StatusCode == http.StatusBadRequest {
		if err = app.errorJSON(w, errors.New("invalid request"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("createUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    jsonResponseBody,
		})
		return
	} else if response.StatusCode != http.StatusCreated {
		if err = app.errorJSON(w, errors.New("error calling user service"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("createUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    jsonResponseBody,
		})
		return
	}

	if err != nil {
		if err = app.errorJSON(w, errors.New("error reading response body"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("createUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    err,
		})
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "success"
	resp.Data = jsonResponseBody

	if err = app.writeJSON(w, http.StatusCreated, resp); err != nil {
		return
	}
}

func (app *Config) getUserRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://user-service/users/%s", id), strings.NewReader(""))
	if err != nil {
		app.sendErrorLog("authenticateUserRequest", errorLog{
			StatusCode: 500,
			Message:    err,
		})
		if err = app.errorJSON(w, err, 500); err != nil {
			return
		}
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		if err = app.errorJSON(w, err, response.StatusCode); err != nil {
			return
		}
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusBadRequest {
		if err = app.errorJSON(w, errors.New("invalid request"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("authenticateUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    err,
		})
		return
	} else if response.StatusCode == http.StatusNotFound {
		app.errorJSON(w, errors.New("not found"), response.StatusCode)
		return
	} else if response.StatusCode != http.StatusOK {
		if err = app.errorJSON(w, errors.New("error calling account service"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("authenticateUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    err,
		})
		return
	}

	maxBytes := 10485376 // 1mgb
	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		if err = app.errorJSON(w, errors.New("error reading response body"), response.StatusCode); err != nil {
			return
		}
		app.sendErrorLog("authenticateUserRequest", errorLog{
			StatusCode: response.StatusCode,
			Message:    err,
		})
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "success"
	resp.Data = jsonResponseBody

	if err = app.writeJSON(w, http.StatusOK, resp); err != nil {
		return
	}
}
