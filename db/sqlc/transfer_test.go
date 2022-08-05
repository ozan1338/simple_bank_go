package db

import (
	"context"
	"testing"

	"github.com/ozan1338/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	arg := CreateTransferParams{
		FromAccountID: 4,
		ToAccountID: 5,
		Amount: util.RandomMoney(),
	}

	_, err := testQueries.CreateTransfer(context.Background(),arg)

	require.NoError(t,err)
}

func TestGetTransfer(t *testing.T) {
	transfer , err := testQueries.GetTransfer(context.Background(), 2)

	require.NoError(t,err)
	require.NotEmpty(t,transfer)
}

func TestUpdateTransfer(t *testing.T) {
	arg := UpdateTransferParams{
		ID: 2,
		Amount: util.RandomMoney(),
	}

	_, err := testQueries.UpdateTransfer(context.Background(), arg)

	require.NoError(t,err)
}

func TestDeleteTransfer(t *testing.T) {
	_, err := testQueries.DeleteTransfer(context.Background(), 4)

	require.NoError(t,err)
}

func TestListTransfer(t *testing.T) {
	transfers, err := testQueries.ListTransfer(context.Background())

	require.NoError(t,err)

	for _, transfer := range transfers {
		require.NotEmpty(t,transfer)
	}
}