package user

import (
	"context"
	"fmt"
	"github.com/bowoBp/LoanFlow/internal/adapter/repository"
	"github.com/bowoBp/LoanFlow/internal/constant"
	"github.com/bowoBp/LoanFlow/internal/domain"
	"github.com/bowoBp/LoanFlow/pkg/environment"
	"github.com/bowoBp/LoanFlow/utils/helper"
	"strconv"
	"time"
)

type (
	UsaCase struct {
		userRepo Repository.UserRepoInterface
		jwt      helper.JwtInterface
		bcrypt   helper.BcryptInterface
		env      environment.Environment
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
			email, password string,
		) (*domians.User, string, string, error)

		RefreshToken(
			ctx context.Context,
			id uint,
			createdAt time.Time,
			token, role, name string,
		) (string, string, error)

		RevokeToken(
			ctx context.Context,
			id uint,
		) error
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
) (*domians.User, string, string, error) {
	user, err := uc.userRepo.CheckEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}

	// verify hashed password
	comparePass := uc.bcrypt.ComparePass([]byte(user.PasswordHash), []byte(password))
	if !comparePass {
		return nil, "", "", err
	}

	//Generate token JWT
	tokenString, errToken := uc.jwt.GenerateToken(user.ID, user.Role, user.Name, user.CreatedAt)
	if errToken != nil {
		return nil, "", "", err
	}

	//Generate refresh token
	refreshToken, err := uc._generateRefreshToken(user.ID, ctx)

	return user, tokenString, refreshToken, nil

}

func (uc UsaCase) RefreshToken(
	ctx context.Context,
	id uint,
	createdAt time.Time,
	token, role, name string,
) (string, string, error) {

	//Generate token JWT
	tokenString, err := uc.jwt.GenerateToken(id, role, name, createdAt)
	if err != nil {
		return "", "", err
	}

	//Generate refresh token
	refreshToken, err := uc._generateRefreshToken(id, ctx)

	return tokenString, refreshToken, nil

}

func (uc UsaCase) _generateRefreshToken(
	id uint,
	ctx context.Context,
) (string, error) {
	// Periksa apakah refresh token sudah ada
	existingRefreshToken, err := uc.userRepo.CheckRefreshToken(ctx, id)
	if err != nil {
		return "", err // Error dari query ke database
	}

	// Generate nilai refresh token baru
	refresh, err := uc.bcrypt.GenerateHashValue(
		uc.env.Get("DEFAULT_SECRET_FORGET_PASSWORD"),
		fmt.Sprintf("%d", id),
		17,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token hash: %w", err)
	}

	// Set waktu kedaluwarsa
	expireAt := time.Now().Add(7 * 24 * time.Hour)

	// Jika refresh token sudah ada, update
	if existingRefreshToken != nil {
		existingRefreshToken.RefreshToken = refresh
		existingRefreshToken.ExpiresAt = expireAt
		existingRefreshToken.UpdatedAt = time.Now()

		// Update refresh token di database
		err = uc.userRepo.UpdateRefreshToken(ctx, existingRefreshToken)
		if err != nil {
			return "", fmt.Errorf("failed to update refresh token: %w", err)
		}

		return existingRefreshToken.RefreshToken, nil
	}

	// Jika tidak ada, buat refresh token baru
	newRefreshToken := &domians.RefreshToken{
		UserID:       id,
		RefreshToken: refresh,
		ExpiresAt:    expireAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Simpan refresh token baru ke database
	_, err = uc.userRepo.StoreRefreshToken(ctx, newRefreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return newRefreshToken.RefreshToken, nil
}

func (uc UsaCase) RevokeToken(
	ctx context.Context,
	id uint,
) error {
	return uc.userRepo.DeleteRefreshTokenByUserID(ctx, id)
}
