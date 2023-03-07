package main

import (
	"database/sql"
	"net/http"

	db "github.com/bugrakocabay/dummy-bank-microservice/transaction-service/db/sqlc"
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

type getTransactionRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTransaction(ctx *gin.Context) {
	var req getTransactionRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transaction, err := server.store.GetTransaction(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transaction)
}

func (server *Server) listTransactions(ctx *gin.Context) {
	transactions, err := server.store.ListTransactions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}
