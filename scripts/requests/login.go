package requests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type LoginPayload struct {
	Action string    `json:"action"`
	Login  LoginData `json:"login"`
}

type LoginData struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Error   bool              `json:"error"`
	Message string            `json:"message"`
	Data    LoginResponseData `json:"data"`
}

type LoginResponseData struct {
	AccessToken string   `json:"access_token"`
	UserData    UserData `json:"user"`
}

func Login(userID string) string {
	requestBody := LoginPayload{
		Action: "login",
		Login: LoginData{
			UserID:   userID,
			Password: "qwerty",
		},
	}

	jsonData, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/handle/users/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("request err: %v", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("response err: %v", err)
	}
	defer response.Body.Close()

	var resp LoginResponse
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Fatalf("failed to decode response body: %v", err)
	}

	return resp.Data.AccessToken
}
