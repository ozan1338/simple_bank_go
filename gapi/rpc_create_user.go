package gapi

import (
	"context"
	"fmt"

	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/pb"
	"github.com/ozan1338/util"
	"github.com/ozan1338/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := ValidateCreateUserRequest(req)

	if violations != nil {
		return nil , invalidArgumentError(violations)
	}
	
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

func ValidateCreateUserRequest(req *pb.CreateUserRequest) (violation []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violation = append(violation, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violation = append(violation, fieldViolation("password", err))
	}

	if err := val.ValidatePassword(req.GetFullName()); err != nil {
		violation = append(violation, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violation = append(violation, fieldViolation("email", err))
	}

	return violation
}

