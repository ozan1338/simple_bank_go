package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashPass, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPass)

	err = CheckPassword(password,hashPass)
	require.NoError(t,err)

	wrongPasswpord := RandomString(6)
	err = CheckPassword(wrongPasswpord,hashPass)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashPass1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPass)
	require.NotEqual(t, hashPass, hashPass1)
}