package main

import (
	db "github.com/bugrakocabay/dummy-bank-microservice/account-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store    db.Store
	router   *gin.Engine
	transfer db.SQLStore
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts/create", server.createAccount)
	router.GET("/accounts/:account_id", server.getAccount)
	router.GET("/accounts/balance/:account_id", server.getAccountBalance)
	router.PUT("/accounts/update", server.updateAccount)
	router.DELETE("/accounts/delete/:account_id", server.deleteAccount)

	router.POST("/transactions/create", server.createTransfer)
	router.GET("/transactions/:transaction_id", server.getTransaction)
	router.GET("/transactions", server.listTransactions)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
