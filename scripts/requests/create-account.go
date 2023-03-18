package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CreateAccountPayload struct {
	Action string            `json:"action"`
	Create CreateAccountData `json:"create"`
}

type CreateAccountData struct {
	Currency string `json:"currency"`
}

type CreateAccountResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    AccountData `json:"data"`
}

type AccountData struct {
	AccountID string    `json:"account_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	Currency  string    `json:"currency"`
	UserID    string    `json:"user_id"`
}

func CreateAccount(accessToken string) string {
	requestBody := CreateAccountPayload{
		Action: "create",
		Create: CreateAccountData{
			Currency: "EUR",
		},
	}

	jsonData, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/handle/accounts", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("request err: %v", err)
	}

	token := fmt.Sprintf("Bearer %s", accessToken)
	request.Header.Set("Authorization", token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("response err: %v", err)
	}
	defer response.Body.Close()

	var resp CreateAccountResponse
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Fatalf("failed to decode response body: %v", err)
	}

	return resp.Data.AccountID
}
