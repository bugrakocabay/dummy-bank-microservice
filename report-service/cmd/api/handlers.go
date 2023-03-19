package main

import (
	"log"
	"time"

	db "github.com/bugrakocabay/dummy-bank-microservice/report-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (server *Server) getDailyReport(ctx *gin.Context) {
	currentTime := time.Now()

	report, err := server.store.GetDailyTransactionReport(ctx, currentTime)
	if err != nil {
		server.sendErrorLog("getDailyReport", Log{
			StatusCode: 500,
			Message:    err,
		})
		return
	}

	resp := db.SaveDailyTransactionReportParams{
		NumTransactions:        int32(report.NumTransactions),
		TotalTransactionAmount: int32(report.TotalTransactionAmount),
		AvgTransactionAmount:   report.AvgTransactionAmount,
		TotalCommission:        report.TotalCommission,
		Day:                    report.Day.Format("02.01.2006"),
	}

	err = server.store.SaveDailyTransactionReport(ctx, resp)
	if err != nil {
		server.sendErrorLog("getDailyReport", Log{
			StatusCode: 500,
			Message:    err,
		})
		return
	}

	log.Println("Saved daily-report cron.")
}
