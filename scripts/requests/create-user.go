package requests

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type CreateUserPayload struct {
	Action string         `json:"action"`
	Create CreateUserData `json:"create"`
}

type CreateUserData struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
}

type GetUserResponse struct {
	Error   bool     `json:"error"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type UserData struct {
	CreatedAt string `json:"created_at"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserID    string `json:"user_id"`
}

func CreateUser() string {
	requestBody := CreateUserPayload{
		Action: "create",
		Create: CreateUserData{
			Firstname: RandomString(5),
			Lastname:  RandomString(6),
			Password:  "qwerty",
		},
	}

	jsonData, _ := json.Marshal(requestBody)
	request, err := http.NewRequest(http.MethodPost, "http://localhost:8080/handle/users", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("request err: %v", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("response err: %v", err)
	}
	defer response.Body.Close()

	var resp GetUserResponse
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Fatalf("failed to decode response body: %v", err)
	}

	return resp.Data.UserID
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
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
