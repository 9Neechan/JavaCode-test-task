package api

import (
	"database/sql"
	"net/http"

	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// gin не может корректно распарсить uuid из uri, пожтому используем string
type getWalletRequest struct {
	WalletUuid string `uri:"id" binding:"required,uuid"`
}

func (server *Server) getWallet(ctx *gin.Context) {
	var req getWalletRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parsedUUID, err := uuid.Parse(req.WalletUuid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	wallet, err := server.store.GetWallet(ctx, parsedUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, wallet.Balance)
}

type UpdateWalletBalanceRequest struct {
	WalletUuid    uuid.UUID `json:"wallet_uuid" binding:"required,uuid"`
	Amount        int64     `json:"amount" binding:"required,gt=0"`
	OperationType string    `json:"operation_type" binding:"required,oneof=DEPOSIT WITHDRAW"`
}

func (server *Server) updateWalletBalance(ctx *gin.Context) {
	var req UpdateWalletBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.TransferTxParams{
		Amount:     req.Amount,
		WalletUuid: req.WalletUuid,
		OperationType: req.OperationType,
	}

	wallet, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}
