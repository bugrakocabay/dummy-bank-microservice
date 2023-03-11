package main

import (
	db "github.com/bugrakocabay/dummy-bank-microservice/user-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/users/create", server.createUser)
	router.GET("/users/:user_id", server.getUser)
	router.POST("/users/authenticate", server.authenticateUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}