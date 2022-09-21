package gapi

import (
	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converUser(user db.CreateUserParams) *pb.User {
	return &pb.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		CreatedAt: timestamppb.Now(),
	}
}