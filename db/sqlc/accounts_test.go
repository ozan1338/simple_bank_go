package db

import (
	"context"
	"testing"

	"github.com/ozan1338/util"
	"github.com/stretchr/testify/require"
)



func TestCreateAccount(t *testing.T) {
	// test := util.RandomOwner()

	user := CreateRandomUser(t)
	arg := CreateAccountParams{
		Owner: user.Username,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	_, err := testQueries.CreateAccount(context.Background(),arg)

	// if err == nil {
	// 	log.Fatal(id)
	// }

	require.NoError(t, err)

}

func TestGetAccount(t *testing.T) {
	account,err := testQueries.GetAccount(context.Background(), 1)
	
	// if err == nil {
	// 	// fmt.Print(account)
	// 	log.Fatal(account)
	// }
	require.NoError(t,err)
	require.NotEmpty(t, account)

	// fmt.Print(account)
}

func TestUpdatedAccount(t *testing.T) {
	arg := UpdateAccountParams{
		ID: 2,
		Balance: util.RandomMoney(),
	}

	_, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t,err)
	// fmt.Print(account)
}

func TestDeleteAccount(t * testing.T) {
	_,err := testQueries.DeleteAccount(context.Background(), 12)

	require.NoError(t,err)
	// require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		arg := CreateAccountParams{
			Owner: util.RandomOwner(),
			Balance: util.RandomMoney(),
			Currency: util.RandomCurrency(),
		}
	
		_, err := testQueries.CreateAccount(context.Background(),arg)

		require.NoError(t,err)
	}

	arg := ListAccountParams{
		Limit: 5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(),arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t,account)
	}
}