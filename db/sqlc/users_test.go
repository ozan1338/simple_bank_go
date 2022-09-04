package db

import (
	"context"
	"testing"

	"github.com/ozan1338/util"
	"github.com/stretchr/testify/require"
)


func CreateRandomUser(t *testing.T) CreateUserParams {
	hashedPassword, err := util.HashPassword(util.RandomString(6))

	require.NoError(t, err)

	arg := CreateUserParams{
		Username: util.RandomOwner(),
		Password: hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	_, _ = testQueries.CreateUser(context.Background(), arg)

	// require.NoError(t, err)

	return arg
}


func TestCreateUser(t *testing.T) {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	arg := CreateUserParams{
		Username: util.RandomOwner(),
		Password: hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	_, err = testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
}
func TestGetUser(t *testing.T)  {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	arg := CreateUserParams{
		Username: util.RandomOwner(),
		Password: hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	_,err = testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)

	username1 := arg.Username

	user, err := testQueries.GetUser(context.Background(), username1)

	require.NoError(t, err)
	require.Equal(t, username1, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Password, user.Password)
}