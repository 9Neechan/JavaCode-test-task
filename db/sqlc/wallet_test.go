package db

import (
	"context"
	"testing"
	"time"

	"github.com/9Neechan/JavaCode-test-task/util"

	"github.com/stretchr/testify/require"
)

// createRandomWallet создает случайный кошелек для тестирования.
func createRandomWallet(t *testing.T) Wallet {
	balance := util.RandomMoney()

	wallet, err := testQueries.CreateWallet(context.Background(), balance)
	require.NoError(t, err)
	require.NotEmpty(t, wallet)

	require.Equal(t, balance, wallet.Balance)

	require.NotZero(t, wallet.WalletUuid)
	require.NotZero(t, wallet.CreatedAt)

	return wallet
}

// TestCreateWallet проверяет создание кошелька.
func TestCreateWallet(t *testing.T) {
	createRandomWallet(t)
}

// TestGetWallet проверяет получение кошелька по его UUID.
func TestGetWallet(t *testing.T) {
	wallet1 := createRandomWallet(t)
	wallet2, err := testQueries.GetWallet(context.Background(), wallet1.WalletUuid)

	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.WalletUuid, wallet2.WalletUuid)
	require.Equal(t, wallet1.Balance, wallet2.Balance)
	// проверить, что две метки времени отличаются не более чем на 1 секунду
	require.WithinDuration(t, wallet1.CreatedAt, wallet2.CreatedAt, time.Second)
}

// TestUpdateWalletBalance проверяет обновление баланса кошелька.
func TestUpdateWalletBalance(t *testing.T) {
	wallet1 := createRandomWallet(t)

	arg := UpdateWalletBalanceParams{
		WalletUuid: wallet1.WalletUuid,
		Amount:   util.RandomMoney(),
	}

	wallet2, err := testQueries.UpdateWalletBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, wallet2)

	require.Equal(t, wallet1.WalletUuid, wallet2.WalletUuid)
	require.Equal(t, wallet1.Balance+arg.Amount, wallet2.Balance)
	require.WithinDuration(t, wallet1.CreatedAt, wallet2.CreatedAt, time.Second)
}
