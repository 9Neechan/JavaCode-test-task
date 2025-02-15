package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Тест на успешный депозит (DEPOSIT)
// Проверяет, что депозит средств на кошелек происходит успешно и баланс увеличивается на сумму депозита.
func TestTransferTx_Deposit(t *testing.T) {
	wallet := createRandomWallet(t)
	store := NewStore(testDB)

	amount := int64(1000)

	arg := TransferTxParams{
		Amount:        amount,
		WalletUuid:    wallet.WalletUuid,
		OperationType: "DEPOSIT",
	}

	// Выполняем транзакцию
	result, err := store.TransferTx(context.Background(), arg)

	// Проверяем, что ошибок нет
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Проверяем, что баланс увеличился на сумму депозита
	require.Equal(t, arg.WalletUuid, result.WalletUuid)
	require.Equal(t, arg.Amount+wallet.Balance, result.Balance)
}

// Тест на успешное снятие средств (WITHDRAW)
// Проверяет, что снятие средств с кошелька происходит успешно и баланс уменьшается на сумму вывода.
func TestTransferTx_Withdraw(t *testing.T) {
	store := NewStore(testDB)

	wallet := createRandomWallet(t)
	withdrawAmount := int64(10)

	// Устанавливаем начальный баланс
	_, err := store.TransferTx(context.Background(), TransferTxParams{
		Amount:        withdrawAmount,
		WalletUuid:    wallet.WalletUuid,
		OperationType: "DEPOSIT",
	})
	require.NoError(t, err)

	arg := TransferTxParams{
		Amount:        withdrawAmount,
		WalletUuid:    wallet.WalletUuid,
		OperationType: "WITHDRAW",
	}

	// Выполняем транзакцию
	result, err := store.TransferTx(context.Background(), arg)

	// Проверяем, что ошибок нет
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Проверяем, что баланс уменьшился на сумму вывода
	require.Equal(t, wallet.WalletUuid, result.WalletUuid)
	require.Equal(t, wallet.Balance, result.Balance)
}

// Тест на недостаток средств при снятии
// Проверяет, что при снятии средств с кошелька, если средств недостаточно, возвращается ошибка "Недостаточно средств".
func TestTransferTx_InsufficientFunds(t *testing.T) {
	store := NewStore(testDB)

	wallet := createRandomWallet(t)
	withdrawAmount := int64(10000)

	arg := TransferTxParams{
		Amount:        withdrawAmount,
		WalletUuid:    wallet.WalletUuid,
		OperationType: "WITHDRAW",
	}

	// Выполняем транзакцию
	_, err := store.TransferTx(context.Background(), arg)

	// Ожидаем ошибку "Недостаточно средств"
	require.Error(t, err)
	require.EqualError(t, err, "Insufficient funds, Недостаточно средств")
}


func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	wallet := createRandomWallet(t)
	fmt.Println("  >> before:", wallet.Balance)

	// run a concurrent transactions
	n := 2
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	// запускаем n одновременных транзакций перевода
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		go func() {
			// выполняется внутри другой горутины, отличной от TestTransferTx
			// => не можем использовать testify/require
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				Amount:        amount,
				WalletUuid:    wallet.WalletUuid,
				OperationType: "WITHDRAW",
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		expectedBalance := wallet.Balance - int64(i+1)*amount
		require.Equal(t, expectedBalance, result.Balance)
		require.True(t, expectedBalance > 0)
	}

	// check final updated wallet
	updatedWallet, err := store.GetWallet(context.Background(), wallet.WalletUuid)
	require.NoError(t, err)

	fmt.Println("  >> after:", updatedWallet.Balance)
	require.Equal(t, wallet.Balance-int64(n)*amount, updatedWallet.Balance)
}

// проверяем только на deadlock, результат проверили в функции выше
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	wallet := createRandomWallet(t)
	fmt.Println("  >> before:", wallet.Balance)

	n := 10
	amount := int64(1)
	errs := make(chan error)

	// запускаем n одновременных транзакций перевода
	for i := 0; i < n; i++ {
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				Amount:     amount,
				WalletUuid: wallet.WalletUuid,
				OperationType: "WITHDRAW",
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check final updates
	updatedWallet, err := store.GetWallet(context.Background(), wallet.WalletUuid)
	require.NoError(t, err)

	fmt.Println("  >> after:", updatedWallet.Balance)
	require.Equal(t, wallet.Balance-int64(n)*amount, updatedWallet.Balance)
}
