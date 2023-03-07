package main

import (
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

	// Accounts-services
	mux.Post("/handle/accounts", app.HandleAccounts)
	mux.Get("/handle/accounts/{id}", app.HandleAccounts)
	mux.Put("/handle/accounts/update", app.HandleAccounts)
	mux.Delete("/handle/accounts/delete/{id}", app.HandleAccounts)

	// Transactions-services
	mux.Post("/handle/transactions", app.HandleTransactions)
	mux.Get("/handle/transactions/{id}", app.HandleTransactions)
	mux.Get("/handle/transactions", app.HandleTransactions)

	return mux
}
