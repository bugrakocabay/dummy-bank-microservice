package main

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/bugrakocabay/dummy-bank-microservice/account-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransactionRequest struct {
	FromAccountID     string         `json:"from_account_id" binding:"required"`
	ToAccountID       string         `json:"to_account_id" binding:"required"`
	TransactionAmount int32          `json:"transaction_amount" binding:"required"`
	Description       sql.NullString `json:"description"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account1, err := server.store.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if account1.Balance < req.TransactionAmount {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("sender doesn't have enough money")))
		return
	}

	_, err = server.store.GetAccount(ctx, req.ToAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := db.TransferTxParams{
		TransactionID:     server.createUUID(),
		FromAccountID:     req.FromAccountID,
		ToAccountID:       req.ToAccountID,
		TransactionAmount: req.TransactionAmount,
		Description:       req.Description,
	}
	transaction, err := server.store.TransferTx(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, transaction)
}

type getTransactionRequest struct {
	TransactionID string `uri:"transaction_id" binding:"required,min=1"`
}

func (server *Server) getTransaction(ctx *gin.Context) {
	var req getTransactionRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transaction, err := server.store.GetTransaction(ctx, req.TransactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
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
