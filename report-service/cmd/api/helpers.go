package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
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

type JSONPayload struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

type Log struct {
	Message    any `json:"message"`
	StatusCode int `json:"status_code"`
}

func (server *Server) sendErrorLog(name string, payload Log) {
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

	log.Println("error logged successfully")
	return
}

// runDailyCron runs a cron that runs every day at 00:00
func (server *Server) runDailyCron(ctx *gin.Context) {
	var lock sync.Mutex
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	timer := time.NewTimer(next.Sub(now))
	defer timer.Stop()
	for {
		/* run forever */
		select {
		case <-timer.C:
			go func() {
				lock.Lock()
				defer lock.Unlock()
				_, err := http.Get("http://report-service/reports/daily-report")
				if err != nil {
					log.Println("error daily report request: ", err)
				}

				log.Println("ran daily cron: ", now)
			}()
			// reset timer for the next day
			now = time.Now()
			next = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			timer.Reset(next.Sub(now))
		}
	}
}
