package user

import (
	"context"
	"github.com/bowoBp/LoanFlow/internal/adapter/Repository"
	"github.com/bowoBp/LoanFlow/internal/constant"
	"github.com/bowoBp/LoanFlow/internal/domain"
	"github.com/bowoBp/LoanFlow/utils/helper"
	"strconv"
	"time"
)

type (
	UsaCase struct {
		userRepo Repository.UserRepoInterface
		jwt      helper.JwtInterface
		bcrypt   helper.BcryptInterface
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

	emailDuplicate, err := uc.userRepo.CheckEmail(ctx, payload.Email)
	if err != nil {
		return result, err
	}
	if emailDuplicate != nil {
		return result, constant.DuplicateEmail
	}
	hashPass := uc.bcrypt.HasPass(payload.Password)
	user, err := uc.userRepo.StoreUser(
		ctx,
		&domians.User{
			Name:         payload.UserName,
			PasswordHash: hashPass,
			Email:        payload.Email,
			Phone:        payload.Phone,
			Role:         payload.Role,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	result.User = RegisterUser{
		ID:        strconv.Itoa(int(user.ID)),
		UserName:  user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
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
	email, password string,
) (*domians.User, string, error) {
	user, err := uc.userRepo.CheckEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	// verify hashed password
	comparePass := uc.bcrypt.ComparePass([]byte(user.PasswordHash), []byte(password))
	if !comparePass {
		return nil, "", err
	}

	//Generate token JWT
	tokenString, errToken := uc.jwt.GenerateToken(user.ID, user.Role, user.Name, user.CreatedAt)
	if errToken != nil {
		return nil, "", err
	}

	return user, tokenString, nil

}
