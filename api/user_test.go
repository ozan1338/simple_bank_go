package api

import (
	_ "testing"

	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/util"
	_ "github.com/stretchr/testify/require"
)

// func TestCreateUserAPI(t *testing.T) {
// 	user,err := randomUser()

// 	require.NoError(t,err)

// }

func randomUser() (user db.User, err error ){
	hashPass, err := util.HashPassword(util.RandomString(6))
	if err != nil {
		return db.User{},err
	}

	return db.User{
		Username: util.RandomOwner(),
		Password: hashPass,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}, nil
}