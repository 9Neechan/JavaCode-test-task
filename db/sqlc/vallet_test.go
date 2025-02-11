package db

import (
	"context"
	"testing"
	"time"

	"github.com/9Neechan/JavaCode-test-task/util"

	"github.com/stretchr/testify/require"
)

func createRandomVallet(t *testing.T) Vallet {
	balance := util.RandomMoney()

	vallet, err := testQueries.CreateVallet(context.Background(), balance)
	require.NoError(t, err)
	require.NotEmpty(t, vallet)

	require.Equal(t, balance, vallet.Balance)

	require.NotZero(t, vallet.ValletID)
	require.NotZero(t, vallet.CreatedAt)

	return vallet
}

func TestCreateVallet(t *testing.T) {
	createRandomVallet(t)
}

func TestGetVallet(t *testing.T) {
	vallet1 := createRandomVallet(t)
	vallet2, err := testQueries.GetVallet(context.Background(), vallet1.ValletID)

	require.NoError(t, err)
	require.NotEmpty(t, vallet2)

	require.Equal(t, vallet1.ValletID, vallet2.ValletID)
	require.Equal(t, vallet1.Balance, vallet2.Balance)
	// проверить, что две метки времени отличаются не более чем на 1 секунду
	require.WithinDuration(t, vallet1.CreatedAt, vallet2.CreatedAt, time.Second)
}

func TestUpdateValletBalance(t *testing.T) {
	vallet1 := createRandomVallet(t)

	arg := UpdateValletBalanceParams{
		ValletID: vallet1.ValletID,
		Amount:   util.RandomMoney(),
	}

	vallet2, err := testQueries.UpdateValletBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, vallet2)

	require.Equal(t, vallet1.ValletID, vallet2.ValletID)
	require.Equal(t, vallet1.Balance+arg.Amount, vallet2.Balance)
	require.WithinDuration(t, vallet1.CreatedAt, vallet2.CreatedAt, time.Second)
}
