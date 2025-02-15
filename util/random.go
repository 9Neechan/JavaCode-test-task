package util

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

// Инициализация генератора случайных чисел
func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

// RandomInt генерирует случайное целое число в диапазоне от min до max включительно.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomMoney генерирует случайное количество денег в диапазоне от 1000 до 10000 включительно.
func RandomMoney() int64 {
	return RandomInt(1000, 10000)
}

// RandomUUID генерирует случайный UUID.
func RandomUUID() uuid.UUID {
	return uuid.New()
}
