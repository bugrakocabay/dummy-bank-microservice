package main

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/bugrakocabay/dummy-bank-microservice/account-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type accountResponse struct {
	Balance   int32     `json:"balance"`
	Currency  string    `json:"currency"`
	AccountID string    `json:"account_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func newAccountResponse(account db.Account) accountResponse {
	return accountResponse{
		Balance:   account.Balance,
		Currency:  account.Currency,
		AccountID: account.AccountID,
		UserID:    account.UserID,
		CreatedAt: account.CreatedAt,
	}
}

type createAccountRequest struct {
	Currency string `json:"currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := db.CreateAccountParams{
		AccountID: server.createUUID(),
		UserID:    server.createUUID(),
		Currency:  req.Currency,
		Balance:   0,
	}

	account, err := server.store.CreateAccount(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newAccountResponse(account)
	ctx.JSON(http.StatusCreated, resp)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newAccountResponse(account)
	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) getAccountBalance(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccountBalance(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type updateAccountRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int32 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	account, err := server.store.UpdateAccount(ctx, payload)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newAccountResponse(account)
	ctx.JSON(http.StatusOK, resp)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var err error
	_, err = server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, struct{}{})
}
