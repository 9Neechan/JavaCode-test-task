package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	// 1. Читаем баланс из Redis
	cachedBalance, err := redis_cache.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Если баланс найден в кэше, парсим и возвращаем
		balance, parseErr := strconv.ParseInt(cachedBalance, 10, 64)
		if parseErr != nil {
			fmt.Println("Ошибка парсинга баланса из Redis:", parseErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cached balance"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"source": "cache", "balance": balance})
		return
	}

	// 2. Если данных нет в Redis → запрашиваем из БД
	wallet, err := server.store.GetWallet(ctx, parsedUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// 3. Записываем данные в Redis с TTL = 5 секунд
	err = redis_cache.RedisClient.Set(ctx, cacheKey, wallet.Balance, 5*time.Second).Err()
	if err != nil {
		fmt.Println("Ошибка записи в Redis:", err) // Логируем, но продолжаем выполнение
	}

	// 4. Возвращаем данные из БД
	ctx.JSON(http.StatusOK, gin.H{"source": "db", "balance": wallet.Balance})
}
