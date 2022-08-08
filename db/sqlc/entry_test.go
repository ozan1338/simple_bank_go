package db

import (
	"context"
	"testing"

	"github.com/ozan1338/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntries(t *testing.T) {
	arg := CreateEntriesParams{
		AccountID: 4,
		Amount: 85,
	}

	_, err := testQueries.CreateEntries(context.Background(), arg)

	require.NoError(t,err)
}

func TestGetEntries(t *testing.T) {
	entries,err := testQueries.GetEntries(context.Background(), 1)

	require.NoError(t,err)
	require.NotEmpty(t,entries)
}

func TestUpdateEntries(t *testing.T) {
	arg := UpdateEntriesParams{
		Amount: util.RandomMoney(),
		ID: 1,
	}

	_, err := testQueries.UpdateEntries(context.Background(), arg)

	require.NoError(t,err)
}

func TestDeleteEntries(t *testing.T) {
	_, err := testQueries.DeleteEntries(context.Background(), 1)

	require.NoError(t,err)
}

func TestListEntries(t *testing.T) {
	entries, err := testQueries.ListEntries(context.Background())

	require.NoError(t,err)

	for _,entry := range entries {
		require.NotEmpty(t,entry)
	}
}