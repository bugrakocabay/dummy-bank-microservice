package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type dailyReportResponse struct {
	NumTransactions        int64   `json:"num_transactions"`
	AvgTransactionAmount   float64 `json:"avg_transaction_amount"`
	TotalTransactionAmount int64   `json:"total_transaction_amount"`
	TotalCommission        float64 `json:"total_commission"`
	Day                    string  `json:"day"`
}

func (server *Server) getDailyReport(ctx *gin.Context) {
	currentTime := time.Now()

	report, err := server.store.GetDailyTransactionReport(ctx, currentTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := dailyReportResponse{
		NumTransactions:        report.NumTransactions,
		TotalTransactionAmount: report.TotalTransactionAmount,
		AvgTransactionAmount:   report.AvgTransactionAmount,
		TotalCommission:        report.TotalCommission,
		Day:                    report.Day.Format("02.01.2006"),
	}
	ctx.JSON(http.StatusOK, resp)
}
