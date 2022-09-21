package gapi

import (
	"context"
	"fmt"

	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/pb"
	"github.com/ozan1338/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashPass,err := util.HashPassword(req.GetPassword())
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}


	arg := db.CreateUserParams{
		Username: req.GetUsername(),
		Password: hashPass,
		Email: req.GetEmail(),
		FullName: req.GetFullName(),
	}


	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// account := db.
	rsp := &pb.CreateUserResponse{
		User: converUser(arg),
	}
	
	return rsp,nil
}

