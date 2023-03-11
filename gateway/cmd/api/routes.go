package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// routes is responsible for creating new routes for the API and specifying who is allowed to connect with specific
// origins, methods and headers.
func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodGet, http.MethodPatch},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(authenticate)

	// Accounts-services
	mux.Post("/handle/accounts", app.HandleAccounts)
	mux.Get("/handle/accounts/{id}", app.HandleAccounts)
	mux.Put("/handle/accounts/update", app.HandleAccounts)
	mux.Delete("/handle/accounts/delete/{id}", app.HandleAccounts)

	// Transactions-services
	mux.Post("/handle/transactions", app.HandleTransactions)
	mux.Get("/handle/transactions/{id}", app.HandleTransactions)
	mux.Get("/handle/transactions", app.HandleTransactions)

	// Users-services
	mux.Post("/handle/users", app.HandleUsers)
	mux.Get("/handle/users/{user_id}", app.HandleUsers)
	mux.Post("/handle/users/authenticate", app.HandleUsers)

	return mux
}

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the Authentication service to validate the access token
		// If the token is valid, call the next handler
		// If the token is not valid, return an error response
		log.Println("helloww")
		next.ServeHTTP(w, r)
	})
}
