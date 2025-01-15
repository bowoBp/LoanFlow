package user

import (
	"github.com/bowoBp/LoanFlow/internal/constant"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	RequestHandler struct {
		ctrl ControllerInterface
	}
)

func (rh RequestHandler) Register(ctx *gin.Context) {
	var payload = RegisterUser{}
	err := ctx.Bind(&payload)
	role := ctx.GetHeader("X-User-Role")
	payload.Role = role
	if payload.Password == "" || payload.UserName == "" {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorInvalidDataWithMessage(constant.ErrRegister.Error()))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorInvalidDataWithMessage(err.Error()))
		return
	}
	res, err := rh.ctrl.Register(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (rh RequestHandler) GetAll(ctx *gin.Context) {
	res, err := rh.ctrl.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
	}
	ctx.JSON(http.StatusOK, res)

}

func (rh RequestHandler) GetCurrent(ctx *gin.Context) {
	id, _ := ctx.Get("id")
	userName, _ := ctx.Get("userName")
	created, _ := ctx.Get("createdAt")
	res, err := rh.ctrl.GetCurrent(id, userName, created)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(constant.ErrUserNotFound.Error()))
	}
	ctx.JSON(http.StatusOK, res)

}

func (rh RequestHandler) LoginCustomer(ctx *gin.Context) {
	var payload = LoginParam{}
	err := ctx.Bind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorInvalidDataWithMessage(err.Error()))
		return
	}
	res, err := rh.ctrl.Login(ctx, payload.Email, payload.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
	}

	ctx.JSON(http.StatusOK, res)
}
