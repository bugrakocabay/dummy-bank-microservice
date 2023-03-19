package main

import (
	db "github.com/bugrakocabay/dummy-bank-microservice/report-service/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.GET("/reports/daily-report", server.getDailyReport)

	c := cron.New()
	defer c.Stop()

	ctx := gin.Context{}
	c.AddFunc("*/10 * * * * *", func() {
		server.getDailyReport(&ctx)
	})
	c.Start()

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
