package user

import (
	"github.com/bowoBp/LoanFlow/internal/constant"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

func (rh RequestHandler) Login(ctx *gin.Context) {
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

func (rh RequestHandler) Logout(ctx *gin.Context) {
	id, ok := ctx.Get("id")
	if !ok {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorResponseWithMessage("Invalid User ID type"))
	}
	res, err := rh.ctrl.RevokeToken(ctx, id.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
	}
	ctx.JSON(http.StatusOK, res)
}

func (rh RequestHandler) RefreshToken(ctx *gin.Context) {
	// Ambil data dari context
	id, _ := ctx.Get("id")

	idUint, ok := id.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.DefaultErrorResponseWithMessage("Invalid User ID type"))
		return
	}

	userName, exists := ctx.Get("userName")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.DefaultErrorResponseWithMessage("User name not found in context"))
		return
	}
	userNameStr, ok := userName.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.DefaultErrorResponseWithMessage("Invalid UserName type"))
		return
	}

	created, exists := ctx.Get("createdAt")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.DefaultErrorResponseWithMessage("Token creation time not found in context"))
		return
	}
	createdAt, ok := created.(time.Time)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.DefaultErrorResponseWithMessage("Invalid CreatedAt type"))
		return
	}

	role, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.DefaultErrorResponseWithMessage("User role not found in context"))
		return
	}
	roleStr, ok := role.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.DefaultErrorResponseWithMessage("Invalid UserRole type"))
		return
	}

	// Bind payload dari request body
	var payload = RefreshTokenParam{}
	err := ctx.Bind(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.DefaultErrorInvalidDataWithMessage(err.Error()))
		return
	}

	// Panggil controller untuk memproses refresh token
	res, err := rh.ctrl.RefreshToken(
		ctx,
		idUint,
		createdAt,
		payload.RefreshToken,
		roleStr,
		userNameStr,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.DefaultErrorResponseWithMessage(err.Error()))
		return
	}

	// Berikan respons sukses
	ctx.JSON(http.StatusOK, res)
}
