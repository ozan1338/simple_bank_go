package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTX(t *testing.T) {
	store := NewStore(testDb)

	// run n concurrent transfer transaction
	n := 2
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	account1, _ := testQueries.GetAccount(context.Background(), 4)

	account2, _ := testQueries.GetAccount(context.Background(), 2)

	fmt.Println(">> before: ", account1.Balance, account2.Balance)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		FromAccountID := account1.ID
		ToAccountID := account2.ID

		// if i%2 == 1 {
		// 	FromAccountID = account2.ID
		// 	ToAccountID = account1.ID
		// }
		
		go func() {
			fmt.Println("KEPANGGIL")
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParam{
				FromAccountID: FromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check result
	// existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		// require.Equal(t, int64(4), transfer.FromAccountID)

		// fmt.Print(transfer)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)

		// fmt.Print(fromEntry)

		_, err = store.GetEntries(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)

		// fmt.Print(toEntry)

		_, err = store.GetEntries(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)

		// check account balance
		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 :=  toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // amount , 2* amount, 3 * amount ...

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		// require.NotContains(t, existed, k)
		// existed[k] = true

		// TODO: check account balance
		
	}
	//check the final updated balances
	updateAccount1, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccountForUpdate(context.Background(), account2.ID)

	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

}
