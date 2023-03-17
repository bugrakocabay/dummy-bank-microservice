package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type GetAllAccountsPayload struct {
	Action string `json:"action"`
}

type GetAllAccountsResponse struct {
	Error   bool                 `json:"error"`
	Message string               `json:"message"`
	Data    []GetAllAccountsData `json:"data"`
}

type GetAllAccountsData struct {
	AccountID string    `json:"account_id"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	Currency  string    `json:"currency"`
	ID        int       `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id"`
}

func GetAllAccounts(accessToken string) []string {
	requestBody := GetAllAccountsPayload{
		Action: "list",
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

	var resp GetAllAccountsResponse
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Fatalf("failed to decode response body: %v", err)
	}

	var accountIDs []string
	accounts := resp.Data
	for _, i := range accounts {
		accountIDs = append(accountIDs, i.AccountID)
	}

	return accountIDs
}
