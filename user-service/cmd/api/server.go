package main

import (
	"fmt"
	"github.com/bugrakocabay/dummy-bank-microservice/user-service/cmd/token"
	db "github.com/bugrakocabay/dummy-bank-microservice/user-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker("12345678901234567890123456789012") // TODO: move to env variable
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	router.POST("/users/create", server.createUser)
	router.GET("/users/:user_id", server.getUser)
	router.POST("/users/authenticate", server.authenticateUser)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
