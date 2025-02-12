package util

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomMoney() int64 {
	return RandomInt(1000, 10000)
}

func RandomUUID() uuid.UUID {
	return uuid.New()
}
