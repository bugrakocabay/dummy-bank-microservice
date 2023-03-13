package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

// routes is responsible for creating new routes for the API and specifying who is allowed to connect with specific
// origins, methods and headers.
func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(r chi.Router) {
		r.Use(authenticate)
		r.Mount("/handle", app.handleRouter())
	})

	mux.Post("/handle/users/login", app.HandleUsers)
	mux.Post("/handle/users", app.HandleUsers)

	return mux
}

func (app *Config) handleRouter() http.Handler {
	mux := chi.NewRouter()

	// Accounts-services
	mux.Post("/accounts", app.HandleAccounts)
	mux.Get("/accounts/{account_id}", app.HandleAccounts)
	mux.Put("/accounts/update", app.HandleAccounts)
	mux.Delete("/accounts/delete/{account_id}", app.HandleAccounts)

	// Transactions-services
	mux.Post("/transactions", app.HandleTransactions)
	mux.Get("/transactions/{transaction_id}", app.HandleTransactions)
	mux.Get("/transactions", app.HandleTransactions)

	// Users-services
	mux.Get("/users/{user_id}", app.HandleUsers)

	return mux
}

type authenticateResponse struct {
	Status  string          `json:"status"`
	Payload responsePayload `json:"payload"`
}

type responsePayload struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			// Return an error response if the token is missing
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call the Authentication service to validate the access token
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://user-service/users/authenticate"), strings.NewReader(""))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		request.Header.Set("Authorization", token)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var jsonResponseBody authenticateResponse
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&jsonResponseBody)

		ctx := context.WithValue(r.Context(), "user_id", jsonResponseBody.Payload.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
