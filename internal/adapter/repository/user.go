package Repository

import (
	"context"
	"errors"
	"github.com/bowoBp/LoanFlow/internal/domain"
	"gorm.io/gorm"
)

type (
	UserRepo struct {
		db *gorm.DB
	}

	UserRepoInterface interface {
		StoreUser(
			ctx context.Context,
			cust *domians.User,
		) (*domians.User, error)

		GetAllUser(
			ctx context.Context,
		) ([]domians.User, error)

		CheckEmail(
			ctx context.Context,
			userName string,
		) (*domians.User, error)

		StoreRefreshToken(
			ctx context.Context,
			token *domians.RefreshToken,
		) (*domians.RefreshToken, error)

		CheckRefreshToken(
			ctx context.Context,
			id uint,
		) (*domians.RefreshToken, error)

		UpdateRefreshToken(
			ctx context.Context,
			token *domians.RefreshToken) error

		DeleteRefreshTokenByUserID(
			ctx context.Context,
			userID uint) error
	}
)

func NewUserRepo(db *gorm.DB) UserRepoInterface {
	return UserRepo{db: db}
}

func (repo UserRepo) StoreUser(
	ctx context.Context,
	cust *domians.User,
) (*domians.User, error) {
	err := repo.db.WithContext(ctx).
		Create(&cust).
		Error
	return cust, err
}

func (repo UserRepo) StoreRefreshToken(
	ctx context.Context,
	token *domians.RefreshToken,
) (*domians.RefreshToken, error) {
	err := repo.db.WithContext(ctx).
		Create(&token).
		Error
	return token, err
}

func (repo UserRepo) GetAllUser(
	ctx context.Context,
) ([]domians.User, error) {
	var cust []domians.User
	err := repo.db.WithContext(ctx).Find(&cust).
		Error
	return cust, err
}

func (repo UserRepo) CheckEmail(
	ctx context.Context,
	email string,
) (*domians.User, error) {
	customer := &domians.User{}

	err := repo.db.WithContext(ctx).
		Model(&domians.User{}).
		Where("email = ? ", email).
		First(customer).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return customer, err
}

func (repo UserRepo) CheckRefreshToken(
	ctx context.Context,
	id uint,
) (*domians.RefreshToken, error) {
	token := &domians.RefreshToken{}

	err := repo.db.WithContext(ctx).
		Model(&domians.RefreshToken{}).
		Where("user_id = ? ", id).
		First(token).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return token, err
}
func (repo UserRepo) UpdateRefreshToken(
	ctx context.Context,
	token *domians.RefreshToken) error {
	return repo.db.WithContext(ctx).Save(token).Error
}

func (repo UserRepo) DeleteRefreshTokenByUserID(
	ctx context.Context,
	userID uint) error {
	return repo.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&domians.RefreshToken{}).Error
}
