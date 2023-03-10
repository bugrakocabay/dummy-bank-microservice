package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 10485376 // 1mgb
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

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
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

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	log.Println(err)
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

type JSONPayload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

type errorLog struct {
	Message    any `json:"message"`
	StatusCode int `json:"status_code"`
}

func (app *Config) sendErrorLog(name string, payload errorLog) error {
	arg := JSONPayload{
		Name: name,
		Data: payload,
	}
	jsonData, err := json.Marshal(arg)
	if err != nil {
		log.Println("sendErrorLog error: cant marshal json:", err)
		return err
	}
	request, err := http.NewRequest(http.MethodPost, "http://logger-service/logs/create", bytes.NewBuffer(jsonData))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("sendErrorLog error: cant send response:", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		log.Println("sendErrorLog error: calling auth service:", err)
		return err
	}

	log.Println("logged successfully")
	return nil
}
