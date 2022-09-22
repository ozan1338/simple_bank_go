package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenrRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenrResponse struct {
	AccessToken string `json:"access_token"`
	AccessTokenExpires time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenrRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshPayload, err :=server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// sessionID, err := server.store.GetLastInsertId(ctx)
	// refreshPayloadID := 
	session, err := server.store.GetSession(ctx, refreshPayload.ID.String())

	if session.IsBlocked {
		err := fmt.Errorf("blocked Session")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("Incorrect session token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiredAt) {
		err := fmt.Errorf("token expired")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accesToken, accesPayload,err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenrResponse{
		AccessToken: accesToken,
		AccessTokenExpires: accesPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}