package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only one single JSON value")
	}

	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, name string, status int, data any, headers ...http.Header) error {
	logPayload := Log{
		StatusCode: status,
		Message:    data,
	}
	defer app.sendRequestLog(name, logPayload)

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, name string, err error, status ...int) error {
	logPayload := Log{
		StatusCode: status[0],
		Message:    fmt.Sprintf("%v", err),
	}
	defer app.sendErrorLog(name, logPayload)

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = "fail"
	payload.Data = err.Error()

	return app.writeJSON(w, "error", statusCode, payload)
}

type JSONPayload struct {
	Name string `json:"name"`
	Data Log    `json:"data"`
}

type Log struct {
	Message    any `json:"message"`
	StatusCode int `json:"status_code"`
}

func (app *Config) sendErrorLog(name string, payload Log) {
	arg := JSONPayload{
		Name: name,
		Data: payload,
	}
	jsonData, err := json.Marshal(arg)
	if err != nil {
		log.Println("sendErrorLog error: cant marshal json:", err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, "http://logger-service/logs/create-error", bytes.NewBuffer(jsonData))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("sendErrorLog error: cant send response:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		log.Println("sendErrorLog error: calling auth service:", err)
		return
	}

	log.Println("logged successfully")
	return
}

func (app *Config) sendRequestLog(name string, payload Log) {
	arg := JSONPayload{
		Name: name,
		Data: payload,
	}
	jsonData, err := json.Marshal(arg)
	if err != nil {
		log.Println("sendErrorLog error: cant marshal json:", err)
		return
	}
	request, err := http.NewRequest(http.MethodPost, "http://logger-service/logs/create-request", bytes.NewBuffer(jsonData))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("sendErrorLog error: cant send response:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		log.Println("sendErrorLog error: calling auth service:", err)
		return
	}

	log.Println("logged successfully")
	return
}

// getAccountUserID fetches the user ID of the given account
func getAccountUserID(accountID string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://account-service/accounts/%s", accountID), nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var account accountResponse
	err = json.Unmarshal(body, &account)
	if err != nil {
		return "", err
	}

	// Return the user ID
	return account.UserID, nil
}
