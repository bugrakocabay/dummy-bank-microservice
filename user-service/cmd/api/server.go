package main

import (
	"fmt"
	"github.com/bugrakocabay/dummy-bank-microservice/user-service/cmd/token"
	db "github.com/bugrakocabay/dummy-bank-microservice/user-service/db/sqlc"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store db.Store) (*Server, error) {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Error with loading env: ", err)
	}

	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)
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
	router.POST("/users/login", server.loginUser)
	router.GET("/users/authenticate", server.authenticateUser)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
