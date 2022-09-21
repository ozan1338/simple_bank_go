package gapi

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/pb"
	"github.com/ozan1338/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User Not Found")
		}

		return nil, status.Errorf(codes.Internal, "err: %s", err)
	}

	err = util.CheckPassword(req.GetPassword(), user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "not authorize")
	}

	accesToken, accesPayload, err := server.tokenMaker.CreateToken(req.GetPassword(), server.config.AccessTokenDuration)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "err: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "err : %s", err)
	}

	mtdt := server.extractMetaData(ctx)

	_, err = server.store.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent, //TODO: fill it
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		fmt.Println("YOO")
		return nil, status.Errorf(codes.Internal, "err: %s", err)
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)

	arg := db.CreateUserParams{
		Username: user.Username,
		Password: user.Password,
		FullName: user.FullName,
		Email:    user.Email,
	}

	string_uuid := (session.ID).String()

	rsp := &pb.LoginUserResponse{
		User:                  converUser(arg),
		SessionId:             string_uuid,
		AccessToken:           accesToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  timestamppb.New(accesPayload.ExpiredAt),
		RefreshTokenExpiredAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}