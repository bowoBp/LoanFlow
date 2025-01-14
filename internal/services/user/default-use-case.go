package user

import (
	"context"
	"github.com/bowoBp/LoanFlow/internal/adapter/Repository"
	"github.com/bowoBp/LoanFlow/internal/domain"
	"github.com/bowoBp/LoanFlow/utils/helper"
	"time"
)

type (
	UsaCase struct {
		userRepo Repository.UserRepoInterface
	}

	UsecaseInterface interface {
		RegisterUser(
			ctx context.Context,
			payload RegisterUser,
		) (result UseCaseRegisterResult, err error)

		GetAll(
			ctx context.Context,
		) ([]domians.User, error)

		LoginUser(
			ctx context.Context,
			userName, password string,
		) (*domians.User, string, error)
	}
)

func (uc UsaCase) RegisterUser(
	ctx context.Context,
	payload RegisterUser,
) (result UseCaseRegisterResult, err error) {

	hashPass := helper.HasPass(payload.Password)
	user, err := uc.userRepo.StoreUser(
		ctx,
		&domians.User{
			UserName:  payload.UserName,
			Password:  hashPass,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	result.User = RegisterUser{
		ID:        user.ID,
		UserName:  user.UserName,
		CreatedAt: user.CreatedAt,
	}
	return result, err
}

func (uc UsaCase) GetAll(
	ctx context.Context,
) ([]domians.User, error) {
	return uc.userRepo.GetAllUser(ctx)
}

func (uc UsaCase) LoginUser(
	ctx context.Context,
	userName, password string,
) (*domians.User, string, error) {
	user, err := uc.userRepo.LoginUser(ctx, userName)
	if err != nil {
		return nil, "", err
	}

	// verify hashed password
	comparePass := helper.ComparePass([]byte(user.Password), []byte(password))
	if !comparePass {
		return nil, "", err
	}

	//Generate token JWT
	tokenString, errToken := helper.GenerateToken(user.ID, userName, user.CreatedAt)
	if errToken != nil {
		return nil, "", err
	}

	return user, tokenString, nil

}
