package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/ozan1338/db/sqlc"
)


type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

type accountStruct struct {
	Id int64 `json:"id"`
	Owner string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type errorRes struct {
	Error string
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// log.Fatal("ERROR: ",err)
		fmt.Println("ERROR: ",err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	}

	userExist,err := server.store.UserExist(ctx, arg.Owner)

	if !userExist {
		errorRess := errorRes{
			Error: "User Doesn't Exist",
		}
		ctx.JSON(http.StatusForbidden, errorRess)
		return
	}

	userMoreThanOne, err := server.store.UserMoreThanOne(ctx, arg.Owner)
	
	if userMoreThanOne > 1 {
		errorRess := errorRes{
			Error: "Username Already Exist",
		}
		ctx.JSON(http.StatusForbidden, errorRess)
		return
	}

	_, err = server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	accountId,err := server.store.GetLastInsertId(ctx)
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

	account := accountStruct{
		Id: accountId,
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: int64(0),
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccountById(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account,err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
	return
}

type listAccountRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest 
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccount(ctx, arg)

	// fmt.Println(">>>", req)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(err)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
	return
}