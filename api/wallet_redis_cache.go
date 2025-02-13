package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/9Neechan/JavaCode-test-task/db/sqlc"
	"github.com/9Neechan/JavaCode-test-task/redis_cache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (server *Server) getWalletRedis(ctx *gin.Context) {
	var req getWalletRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parsedUUID, err := uuid.Parse(req.WalletUuid)
	if err != nil || parsedUUID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	cacheKey := fmt.Sprintf("wallet:%s", parsedUUID.String())

	// Проверяем кэш Redis
	cachedBalance, err := redis_cache.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Данные найдены в кэше
		ctx.JSON(http.StatusOK, gin.H{"source": "cache", "balance": cachedBalance})
		return
	}

	// Запрашиваем данные из БД
	wallet, err := server.store.GetWallet(ctx, parsedUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Кэшируем баланс в Redis на 5 sec
	err = redis_cache.RedisClient.Set(ctx, cacheKey, wallet.Balance, 5*time.Second).Err()
	if err != nil {
		fmt.Println("Ошибка записи в Redis:", err) // Логируем, но не прерываем выполнение
	}

	ctx.JSON(http.StatusOK, gin.H{"source": "db", "balance": wallet.Balance})
}

func (server *Server) updateWalletBalanceRedis(ctx *gin.Context) {
	var req UpdateWalletBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	parsedUUID, err := uuid.Parse(req.WalletUuid)
	if err != nil || parsedUUID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	arg := db.TransferTxParams{
		Amount:        req.Amount,
		WalletUuid:    parsedUUID,
		OperationType: req.OperationType,
	}

	// Выполняем обновление баланса в БД
	wallet, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Удаляем кеш, чтобы после обновления API не возвращал устаревшие данные
	cacheKey := fmt.Sprintf("wallet:%s", parsedUUID.String())
	err = redis_cache.RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		fmt.Println("Ошибка удаления кеша в Redis:", err) // Логируем, но не прерываем выполнение
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Balance updated", "wallet": wallet})
}

