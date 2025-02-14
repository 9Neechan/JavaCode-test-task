package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/9Neechan/JavaCode-test-task/redis_cache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) updateWalletBalanceRabbitmq(ctx *gin.Context) {
	var req UpdateWalletBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Отправляем сообщение в RabbitMQ
	err := server.rabbitMQ.PublishMessage("wallet_updates", req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Запрос на обновление баланса отправлен"})
}

// processUpdateWallet преобразует []byte в JSON и вызывает обновление баланса
func (server *Server) processUpdateWallet(msg []byte) {
	var req UpdateWalletBalanceRequest

	// Разбираем JSON из RabbitMQ
	if err := json.Unmarshal(msg, &req); err != nil {
		fmt.Println("Ошибка парсинга JSON из RabbitMQ:", err)
		return
	}

	// Вызываем обработчик обновления баланса
	server.handleWalletUpdate(req)
}

// handleWalletUpdate выполняет обновление баланса в БД
func (server *Server) handleWalletUpdate(req UpdateWalletBalanceRequest) {
	parsedUUID, err := uuid.Parse(req.WalletUuid)
	if err != nil || parsedUUID == uuid.Nil {
		fmt.Println("Ошибка: Некорректный UUID")
		return
	}

	arg := db.TransferTxParams{
		Amount:        req.Amount,
		WalletUuid:    parsedUUID,
		OperationType: req.OperationType,
	}

	result, err := server.store.TransferTx(context.Background(), arg)
	if err != nil {
		fmt.Println("Ошибка обновления баланса в БД:", err)
	}

	//!!!!!!!!!!!!!
	// Update cache
	cacheKey := fmt.Sprintf("wallet:%s", parsedUUID.String())
	err = redis_cache.RedisClient.Set(context.Background(), cacheKey, result.Balance, 5*time.Second).Err()
	if err != nil {
		fmt.Println("Ошибка записи в Redis:", err) // Логируем, но не прерываем выполнение
	}
}
