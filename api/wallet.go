package api

import (
	"database/sql"
	"net/http"

	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type getWalletRequest struct {
	WalletID int64 `uri:"wallet_id" binding:"min=1"`
}

func (server *Server) getWallet(ctx *gin.Context) {
	var req getWalletRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	wallet, err := server.store.GetWallet(ctx, req.WalletID)
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
	Amount        int64  `json:"amount"`
	WalletID      int64  `json:"wallet_id"`
	OperationType string `json:"operation_type"`
}

func (server *Server) updateWalletBalance(ctx *gin.Context) {
	var req UpdateWalletBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.OperationType == "WITHDRAW" {
		req.Amount = -req.Amount
	}

	arg := db.UpdateWalletBalanceParams{
		Amount:   req.Amount,
		WalletID: req.WalletID,
	}

	wallet, err := server.store.UpdateWalletBalance(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}
