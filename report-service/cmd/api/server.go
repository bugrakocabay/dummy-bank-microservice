package main

import (
	db "github.com/bugrakocabay/dummy-bank-microservice/report-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	ctx := gin.Context{}
	go server.runDailyCron(&ctx)
	router := gin.Default()

	router.GET("/reports/daily-report", server.getDailyReport)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
