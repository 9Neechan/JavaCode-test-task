package redis_cache

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Глобальная переменная для клиента Redis
var RedisClient *redis.Client
var ctx = context.Background()

// InitRedis инициализирует клиент Redis
func InitRedis(redisAddr string, redisPassword string, redisDB int) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // Адрес сервера Redis
		Password: redisPassword, // Пароль (если есть)
		DB:       redisDB,       // Используемая база данных
	})

	// Проверяем соединение
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	fmt.Println("Подключено к Redis!")
}
