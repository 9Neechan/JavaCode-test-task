package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// input params of the transfer transaction
type TransferTxParams struct {
	Amount        int64     `json:"amount"`
	WalletUuid    uuid.UUID `json:"wallet_uuid"`
	OperationType string    `json:"operation_type"`
}

type TransferTxResult struct {
	Balance    int64     `json:"balance"`
	WalletUuid uuid.UUID `json:"wallet_uuid"`
}

var txKey = struct{}{}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)

		params := UpdateWalletBalanceParams{
			Amount:     arg.Amount,
			WalletUuid: arg.WalletUuid,
		}

		if arg.OperationType == "WITHDRAW" {
			fmt.Println(txName, "getting wallet from db")
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

		fmt.Println("00000", wallet_result)

		result.Balance = wallet_result.Balance
		result.WalletUuid = wallet_result.WalletUuid

		return err
	})

	return result, err
}
