package db

import (
	"context"
	"database/sql"
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

func TestUpdateUser(t *testing.T) {
	oldUser := CreateRandomUser(t)

	newFullName := util.RandomOwner()
	_, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid: true,
		},
	})

	require.NoError(t, err)

	newUser,err := testQueries.GetUser(context.Background(), oldUser.Username)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, newFullName, newUser.FullName)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.Password, newUser.Password)


}