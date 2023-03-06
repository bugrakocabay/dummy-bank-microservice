package main

import (
	"database/sql"
	db "github.com/bugrakocabay/dummy-bank-microservice/transaction-service/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTransactionRequest struct {
	FromAccountID     string         `json:"from_account_id" binding:"required"`
	ToAccountID       string         `json:"to_account_id" binding:"required"`
	TransactionAmount int32          `json:"transaction_amount" binding:"required"`
	Description       sql.NullString `json:"description"`
}

func (server *Server) createTransaction(ctx *gin.Context) {
	var req createTransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := db.CreateTransactionParams{
		TransactionID:     server.createUUID(),
		FromAccountID:     req.FromAccountID,
		ToAccountID:       req.ToAccountID,
		TransactionAmount: req.TransactionAmount,
		Description:       req.Description,
	}

	transaction, err := server.store.CreateTransaction(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transaction)
}
