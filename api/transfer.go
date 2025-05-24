package api

import (
	"database/sql"
	"fmt"
	"net/http"
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR VND"`
}

func (server *Server) createTransfer(context *gin.Context) {
	var request transferRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	if !server.validAccount(context, request.FromAccountID, request.Currency) || !server.validAccount(context, request.ToAccountId, request.Currency) {
		return
	}

	arg := db.CreateTransferParams{
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountId,
		Amount:        request.Amount,
	}

	result, err := server.store.TransferTx(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "Transfer created!",
		"data":    result,
	})

}

func (server *Server) validAccount(context *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(context, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errResponse(err))
		} else {
			context.JSON(http.StatusInternalServerError, errResponse(err))
		}
		return false
	} else if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency miss match: %s vs %s", accountId, account.Currency, currency)
		context.JSON(http.StatusNotFound, errResponse(err))
		return false
	}
	return true
}
