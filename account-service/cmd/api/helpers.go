package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

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

type JSONPayload struct {
	Name string `json:"name"`
	Data Log    `json:"data"`
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

	log.Println("logged successfully")
	return
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}
