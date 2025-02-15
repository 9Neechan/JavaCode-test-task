package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	// TransferTx выполняет транзакцию перевода с заданными параметрами и возвращает результат или ошибку
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// NewStore создает новый экземпляр Store с использованием переданной базы данных
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// SQLStore предоставляет все функции для выполнения запросов и транзакций в базе данных
// Использует композицию для расширения функциональности Queries
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewSQLStore создает новый экземпляр SQLStore с использованием переданной базы данных
func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx выполняет функцию в рамках транзакции базы данных
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}