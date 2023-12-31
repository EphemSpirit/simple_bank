package api

import (
	"database/sql"
	"net/http"

	db "github.com/EphemSpirit/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err!= nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: 0,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) GetAccount(ctx *gin.Context) {
	 var req getAccountRequest

	 if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	 }

	 account, err := s.store.GetAccount(ctx, int64(req.ID))
	 if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	 }

	 ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) GetAccounts(ctx *gin.Context) {
	 var req listAccountRequest

	 if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	 }

	 arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.PageID -1) * req.PageSize,
	 }

	 accounts, err := s.store.ListAccounts(ctx, arg)
	 if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	 }

	 ctx.JSON(http.StatusOK, accounts)
}