package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// TransferTxParams содержит параметры для транзакции перевода
type TransferTxParams struct {
	Amount        int64     `json:"amount"`         // Сумма перевода
	WalletUuid    uuid.UUID `json:"wallet_uuid"`    // UUID кошелька
	OperationType string    `json:"operation_type"` // Тип операции (DEPOSIT или WITHDRAW)
}

// TransferTxResult содержит результат транзакции перевода
type TransferTxResult struct {
	Balance    int64     `json:"balance"`     // Баланс кошелька после транзакции
	WalletUuid uuid.UUID `json:"wallet_uuid"` // UUID кошелька
}

// txKey используется для хранения имени транзакции в контексте
var txKey = struct{}{}

// TransferTx выполняет транзакцию перевода с заданными параметрами и возвращает результат или ошибку
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	// Выполняем транзакцию в рамках функции execTx
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey) // Получаем имя транзакции из контекста

		// Создаем параметры для обновления баланса кошелька
		params := UpdateWalletBalanceParams{
			Amount:     arg.Amount,
			WalletUuid: arg.WalletUuid,
		}

		if arg.OperationType == "WITHDRAW" {
			fmt.Println(txName, "getting wallet from db") // Получаем кошелек из базы данных
			wallet, err := q.GetWallet(ctx, params.WalletUuid)
			if err != nil {
				return err
			}

			if wallet.Balance < params.Amount {
				return fmt.Errorf("Insufficient funds, Недостаточно средств")
			}
			params.Amount = -params.Amount
		}

		fmt.Println(txName, "updating wallet balance")
		wallet_result, err := q.UpdateWalletBalance(ctx, params)
		if err != nil {
			return err
		}

		result.Balance = wallet_result.Balance       // Устанавливаем баланс в результат
		result.WalletUuid = wallet_result.WalletUuid // Устанавливаем UUID кошелька в результат

		return err
	})

	return result, err
}
