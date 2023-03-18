package main

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/bugrakocabay/dummy-bank-microservice/account-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type accountResponse struct {
	Balance   float64   `json:"balance"`
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
	UserID   string `json:"user_id"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := db.CreateAccountParams{
		AccountID: server.createUUID(),
		UserID:    req.UserID,
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

type addBalanceRequest struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
}

func (server *Server) addAccountBalance(ctx *gin.Context) {
	var req addBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := db.AddAccountBalanceParams{
		AccountID: req.AccountID,
		Amount:    req.Amount,
	}

	account, err := server.store.AddAccountBalance(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newAccountResponse(account)
	ctx.JSON(http.StatusCreated, resp)
}

type getAccountRequest struct {
	AccountID string `uri:"account_id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.AccountID)
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

	account, err := server.store.GetAccountBalance(ctx, req.AccountID)
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
	AccountID string  `json:"account_id" binding:"required"`
	Balance   float64 `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := db.UpdateAccountParams{
		AccountID: req.AccountID,
		Balance:   req.Balance,
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
	AccountID string `uri:"account_id" binding:"required"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.AccountID)
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

func (server *Server) listAccounts(ctx *gin.Context) {
	accounts, err := server.store.ListAccounts(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
