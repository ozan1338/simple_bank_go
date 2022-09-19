package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/ozan1338/db/sqlc"
	"github.com/ozan1338/util"
)

type createUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Fullname string `json:"full_name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type createUserRes struct {
	Username         string       `json:"username"`
	Email            string       `json:"email"`
	FullName         string       `json:"full_name"`
	AccessTokenExpires time.Time `json:"access_token_expired"`
}

func newUserRes (user db.User) createUserRes {
	return createUserRes{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	log.Fatal("AYANAON")
	// 	fmt.Println("ERROR: ", err)
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }

	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		fmt.Println("ERROR !")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPass,err := util.HashPassword(req.Password)
	if err != nil {
		fmt.Println("ERROR: ", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// fmt.Println(req)

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: hashPass,
		Email: req.Email,
		FullName: req.Fullname,
	}


	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		// log.Fatal(">>>>",err)
		// mySqlErr, ok := err.(*mysql.MySQLError); 
		// fmt.Println(mySqlErr.Error())
		// if ok {
		// 	switch mySqlErr.Error() {
		// 	case "unique_violation":
		// 		ctx.JSON(http.StatusForbidden, errorResponse(err))
		// 		return
		// 	}
		// }
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// userMoreThanOne, err := server.store.UserMoreThanOne(ctx, arg.Username)
	
	// if userMoreThanOne > 1 {
	// 	errorRess := errorRes{
	// 		Error: "Username Already Exist",
	// 	}
	// 	ctx.JSON(http.StatusForbidden, errorRess)
	// 	return
	// }

	// accountId, err := server.store.GetLastInsertId(ctx)
	if err != nil {
		//HOW TO HANDLE DATABASE ERROR
		// if mySqlErr, ok := err.(*mysql.MySQLError); ok {
		// 	switch mySqlErr.Error() {
		// 	case "foreign_key_violation", "unique_violation":
		// 		ctx.JSON(http.StatusForbidden, errorResponse(err))
		// 		return
		// 	}
		// }
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// account := db.
	rsp := newUserRes(db.User{
		Username: arg.Username,
		Password: arg.Password,
		Email: arg.Email,
		FullName: arg.FullName,
	})

	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	sessionID uuid.UUID `json:"session_id"`
	AccessToken string `json:"access_token"`
	AccessTokenExpires time.Time `json:"access_token_expires_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpires time.Time `json:"refresh_token_expires_at"`
	User createUserRes `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("User Not Found")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accesToken, accesPayload,err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload,err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_,err = server.store.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID: refreshPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: ctx.Request.UserAgent(),//TODO: fill it
		ClientIp: ctx.ClientIP(),
		IsBlocked: false,
		ExpiredAt: refreshPayload.ExpiredAt,
	})
	if err != nil {
		fmt.Println("YOO")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// sessionID, err := server.store.GetLastInsertId(ctx)
	session, err := server.store.GetSession(ctx, refreshPayload.ID)

	rsp := loginUserResponse{
		sessionID: session.ID,
		AccessToken: accesToken,
		AccessTokenExpires: accesPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpires: refreshPayload.ExpiredAt,
		User: newUserRes(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}