package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AddBalancePayload struct {
	Action  string            `json:"action"`
	Balance CreateBalanceData `json:"balance"`
}

type CreateBalanceData struct {
	AccountID string `json:"account_id"`
	Amount    int32  `json:"amount"`
}

type AddBalanceResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    AccountData `json:"data"`
}

func AddBalance(accessToken, accountID string) error {
	requestBody := AddBalancePayload{
		Action: "balance",
		Balance: CreateBalanceData{
			AccountID: accountID,
			Amount:    1000,
		},
	}

	jsonData, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/handle/accounts/add-balance", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("request err: %v", err)
		return err
	}

	token := fmt.Sprintf("Bearer %s", accessToken)
	request.Header.Set("Authorization", token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("response err: %v", err)
		return err
	}
	defer response.Body.Close()

	var resp CreateAccountResponse
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Fatalf("failed to decode response body: %v", err)
		return err
	}

	return nil
}
