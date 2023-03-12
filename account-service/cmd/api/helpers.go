package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (server *Server) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
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

func (server *Server) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
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

func (server *Server) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	log.Println(err)

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return server.writeJSON(w, statusCode, payload)
}

func (server *Server) createUUID() string {
	// Generate a new UUID
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err)
	}

	// Set the UUID version and variant bits
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0xbf) | 0x80

	// Convert the UUID to a string format
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
