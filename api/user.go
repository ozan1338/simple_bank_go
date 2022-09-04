package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

	userMoreThanOne, err := server.store.UserMoreThanOne(ctx, arg.Username)
	
	if userMoreThanOne > 1 {
		errorRess := errorRes{
			Error: "Username Already Exist",
		}
		ctx.JSON(http.StatusForbidden, errorRess)
		return
	}

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
	rsp := createUserRes{
		Username: arg.Username,
		Email: arg.Email,
		FullName: arg.FullName,
	}

	ctx.JSON(http.StatusOK, rsp)
}