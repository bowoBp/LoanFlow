package user

import (
	"context"
	"fmt"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"log"
	"strconv"
	"time"
)

type (
	Controller struct {
		Uc UsecaseInterface
	}

	ControllerInterface interface {
		Register(
			ctx context.Context,
			payload RegisterUser,
		) (*dto.Response, error)

		GetAll(
			ctx context.Context,
		) (*dto.Response, error)

		GetCurrent(
			id, userName, created interface{},
		) (*dto.Response, error)

		Login(
			ctx context.Context,
			userName, password string,
		) (SuccessLoginUser, error)
	}
)

func (ctrl Controller) Register(
	ctx context.Context,
	payload RegisterUser,
) (*dto.Response, error) {
	start := time.Now()
	result, err := ctrl.Uc.RegisterUser(ctx, payload)
	if err != nil {
		return nil, err
	}
	return dto.NewSuccessResponse(
		result,
		"Register is success",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil
}

func (ctrl Controller) GetAll(
	ctx context.Context,
) (*dto.Response, error) {
	start := time.Now()
	var res = make([]Users, 0)
	users, err := ctrl.Uc.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, user := range users {
		res = append(res, Users{
			ID:        strconv.Itoa(int(user.ID)),
			UserName:  user.Name,
			CreatedAt: user.CreatedAt,
		})
	}
	return dto.NewSuccessResponse(
		res,
		"success get all ",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil

}

func (ctrl Controller) GetCurrent(
	id, userName, created interface{},
) (*dto.Response, error) {
	start := time.Now()
	userId, ok := id.(string)
	if !ok {
		return nil, fmt.Errorf("ID must be a string")
	}

	userNm, ok := userName.(string)
	if !ok {
		return nil, fmt.Errorf("userName must be a string")
	}

	createdUsr, ok := created.(time.Time)
	if !ok {
		return nil, fmt.Errorf("created must be a time.Time")
	}
	var res = Users{
		ID:        userId,
		UserName:  userNm,
		CreatedAt: createdUsr,
	}

	return dto.NewSuccessResponse(
		res,
		"success current user",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil

}

func (ctrl Controller) Login(
	ctx context.Context,
	email, password string,
) (SuccessLoginUser, error) {
	user, tokenString, refreshToken, err := ctrl.Uc.LoginUser(ctx, email, password)
	if err != nil {
		return SuccessLoginUser{}, err
	}
	response := SuccessLoginUser{
		Response: dto.ResponseMeta{
			Success:      true,
			MessageTitle: "login successful",
			Message:      "success",
			ResponseTime: "",
		},
		UserName:     user.Name,
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
	}

	return response, nil
}
