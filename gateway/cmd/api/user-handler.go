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

const (
	userServiceURL = "http://user-service"
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
		if err = app.errorJSON(w, "HandleUsers", err); err != nil {
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
		if err = app.errorJSON(w, "HandleUsers", errors.New(fmt.Sprintf("unknown action type: %s", requestPayload.Action))); err != nil {
			return
		}
		return
	}
}

type LoginUserPayload struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (app *Config) loginUserRequest(w http.ResponseWriter, payload LoginUserPayload) error {
	jsonData, _ := json.Marshal(payload)

	reqURL := fmt.Sprintf("%s/users/login", userServiceURL)
	request, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return app.errorJSON(w, "loginUserRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "loginUserRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "loginUserRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "loginUserRequest", response.StatusCode, resp)
}

type CreateUserPayload struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (app *Config) createUserRequest(w http.ResponseWriter, payload CreateUserPayload) error {
	jsonData, _ := json.Marshal(payload)

	reqURL := fmt.Sprintf("%s/users/create", userServiceURL)
	request, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return app.errorJSON(w, "createUserRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "createUserRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "createUserRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "createUserRequest", response.StatusCode, resp)
}

func (app *Config) getUserRequest(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "user_id")
	idInHeader := r.Context().Value("user_id")
	if id != idInHeader {
		return app.errorJSON(w, "getUserRequest", errors.New("this is not yours"), 403)
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://user-service/users/%s", id), strings.NewReader(""))
	if err != nil {
		return app.errorJSON(w, "getUserRequest", err, 500)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return app.errorJSON(w, "getUserRequest", err, response.StatusCode)
	}
	defer response.Body.Close()

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))

	var jsonResponseBody any
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&jsonResponseBody)
	if err != nil {
		return app.errorJSON(w, "getUserRequest", errors.New("error reading response body"), response.StatusCode)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Data = jsonResponseBody

	if response.StatusCode != http.StatusOK {
		resp.Message = "fail"
	} else {
		resp.Message = "success"
	}

	return app.writeJSON(w, "getUserRequest", response.StatusCode, resp)
}
