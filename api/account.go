package api

import (
	"net/http"
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR VND"`
}

func (server *Server) createAccount(context *gin.Context) {
	var req createAccountRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	account, err := server.store.CreateAccount(context, db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"data": account,
	})
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(context *gin.Context) {
	var req getAccountRequest
	if err := context.ShouldBindUri(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := server.store.GetAccount(context, req.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": account,
	})
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(context *gin.Context) {
	var req listAccountRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	accounts, err := server.store.ListAccount(context, db.ListAccountParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * int32(req.PageSize),
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": accounts,
	})
}
