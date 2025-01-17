package loan

import (
	"github.com/bowoBp/LoanFlow/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type (
	RequestHandler struct {
		ctrl ControllerInterface
	}
)

func (rh RequestHandler) CreateLoan(ctx *gin.Context) {
	var payload = CreateLoanRequest{}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorResponseWithMessage(err.Error()))
		return
	}
	id, _ := ctx.Get("id")
	payload.ID = id.(uint)
	res, err := rh.ctrl.CreateLoan(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (rh RequestHandler) ApproveLoan(ctx *gin.Context) {
	var payload = ApproveLoanRequest{}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorResponseWithMessage(err.Error()))
	}
	id, _ := ctx.Get("id")
	loanIdParam := ctx.Param("loanId")
	loanID, err := strconv.Atoi(loanIdParam)
	res, err := rh.ctrl.ApproveLoan(
		ctx,
		uint(loanID),
		id.(uint),
		payload,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
	}
	ctx.JSON(http.StatusOK, res)
}

func (rh RequestHandler) StoreInvest(ctx *gin.Context) {
	var payload = InvestLoanRequest{}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorResponseWithMessage(err.Error()))
	}
	id, _ := ctx.Get("id")
	loanIdParam := ctx.Param("loanId")
	loanID, err := strconv.Atoi(loanIdParam)
	res, err := rh.ctrl.StoreInvest(
		ctx,
		uint(loanID),
		id.(uint),
		payload,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
	}
	ctx.JSON(http.StatusOK, res)
}

func (rh RequestHandler) DisburseLoan(ctx *gin.Context) {
	var payload = DisburseLoanRequest{}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorResponseWithMessage(err.Error()))
	}
	id, _ := ctx.Get("id")
	loanIdParam := ctx.Param("loanId")
	loanID, err := strconv.Atoi(loanIdParam)
	res, err := rh.ctrl.DisburseLoan(
		ctx,
		uint(loanID),
		id.(uint),
		payload,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
	}
	ctx.JSON(http.StatusOK, res)
}

func (rh RequestHandler) GetLoan(ctx *gin.Context) {
	loanIDParam := ctx.Param("loanId")
	loanID, err := strconv.ParseUint(loanIDParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorResponseWithMessage(err.Error()))
		return
	}
	res, err := rh.ctrl.GetLoan(ctx, uint(loanID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (rh RequestHandler) GetLoans(ctx *gin.Context) {
	query, err := rh.getLoansQuery(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
		return
	}
	res, err := rh.ctrl.GetLoans(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h RequestHandler) getLoansQuery(c *gin.Context) (dto.GetListQuery, error) {
	var query dto.GetListQuery

	if c.Query("perPage") != "" {
		perPage, err := strconv.ParseInt(c.Query("perPage"), 10, 64)
		if err != nil {

			return dto.GetListQuery{}, err
		}
		query.PerPage = int(perPage)
	}

	if c.Query("page") != "" {
		page, err := strconv.ParseInt(c.Query("page"), 10, 64)
		if err != nil {

			return dto.GetListQuery{}, err
		}
		query.Page = int(page)
	}

	query.Search = c.Query("search")

	return query, nil
}
