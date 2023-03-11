package main

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/bugrakocabay/dummy-bank-microservice/user-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type userResponse struct {
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(account db.User) userResponse {
	return userResponse{
		Firstname: account.Firstname,
		Lastname:  account.Lastname,
		UserID:    account.UserID,
		CreatedAt: account.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	payload := db.CreateUserParams{
		UserID:    server.createUUID(),
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Password:  hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)
	ctx.JSON(http.StatusCreated, resp)
}

type getUserRequest struct {
	UserID string `uri:"user_id" binding:"required"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}

type authenticateUserRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (server *Server) authenticateUser(ctx *gin.Context) {
	var req authenticateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = CheckPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}
