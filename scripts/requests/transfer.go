package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CreateTransactionPayload struct {
	Action string                `json:"action"`
	Create CreateTransactionData `json:"create"`
}

type CreateTransactionData struct {
	FromAccountID     string `json:"from_account_id"`
	ToAccountID       string `json:"to_account_id"`
	TransactionAmount int32  `json:"transaction_amount"`
}

type CreateTransactionResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func CreateTransaction(accessToken, fromAccountID, toAccountID string) error {
	requestBody := CreateTransactionPayload{
		Action: "create",
		Create: CreateTransactionData{
			FromAccountID:     fromAccountID,
			ToAccountID:       toAccountID,
			TransactionAmount: int32(RandomInt(1, 999)),
		},
	}

	jsonData, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/handle/transactions", bytes.NewBuffer(jsonData))
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

	return nil
}
