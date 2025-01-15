package Repository

import (
	"context"
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

	return customer, err
}
