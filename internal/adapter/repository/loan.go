package Repository

import (
	"context"
	"errors"
	domians "github.com/bowoBp/LoanFlow/internal/domain"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	LoanRepo struct {
		db *gorm.DB
	}

	LoanRepoInterface interface {
		CreateLoan(
			ctx context.Context,
			loan *domians.Loan,
		) (*domians.Loan, error)
		GetLoanByID(
			ctx context.Context,
			loanID uint,
		) (*domians.Loan, error)
		CreateLoanState(
			ctx context.Context,
			loanState *domians.LoanStateHistory,
		) error
		ApproveDetail(
			ctx context.Context,
			detail *domians.LoanApprovalDetail,
		) error
		DisburseDetail(
			ctx context.Context,
			disbursed *domians.LoanDisbursementDetail,
		) error
		UpdateLoan(
			ctx context.Context,
			loan *domians.Loan,
			updateData map[string]any,
		) error
		InvestLoan(
			ctx context.Context,
			investLoan *domians.LoanInvestor,
		) error
		GetLoans(
			ctx context.Context,
			query dto.GetListQuery,
		) ([]domians.Loan, int64, error)
	}
)

func NewLoanRepo(db *gorm.DB) LoanRepoInterface {
	return &LoanRepo{
		db: db,
	}
}

func (repo LoanRepo) CreateLoan(
	ctx context.Context,
	loan *domians.Loan,
) (*domians.Loan, error) {
	err := repo.db.WithContext(ctx).
		Create(loan).
		Error
	return loan, err
}

func (repo LoanRepo) CreateLoanState(
	ctx context.Context,
	loanState *domians.LoanStateHistory,
) error {
	return repo.db.WithContext(ctx).Create(loanState).Error
}

func (repo LoanRepo) GetLoanByID(
	ctx context.Context,
	loanID uint,
) (*domians.Loan, error) {
	var loan domians.Loan
	if err := repo.db.WithContext(ctx).
		First(&loan, loanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &loan, nil
}

func (repo LoanRepo) UpdateLoan(
	ctx context.Context,
	loan *domians.Loan,
	updateData map[string]any,
) error {
	return repo.db.WithContext(ctx).
		Omit(clause.Associations).
		Model(&loan).
		Updates(updateData).
		Error
}

func (repo LoanRepo) ApproveDetail(
	ctx context.Context,
	detail *domians.LoanApprovalDetail,
) error {
	return repo.db.WithContext(ctx).Create(detail).Error
}

func (repo LoanRepo) DisburseDetail(
	ctx context.Context,
	disbursed *domians.LoanDisbursementDetail,
) error {
	return repo.db.WithContext(ctx).Create(disbursed).Error
}

func (repo LoanRepo) InvestLoan(
	ctx context.Context,
	investLoan *domians.LoanInvestor,
) error {
	return repo.db.WithContext(ctx).Create(investLoan).Error
}

func (repo LoanRepo) GetLoans(
	ctx context.Context,
	query dto.GetListQuery,
) ([]domians.Loan, int64, error) {
	var (
		loans = make([]domians.Loan, 0)
		count int64
	)

	db := repo.db.WithContext(ctx).
		Model(&domians.Loan{})
	if query.Search != "" {
		db.Where("state like ?", "%"+query.Search+"%")
	}
	dbCount := db
	err := dbCount.Count(&count).Error

	offSet := query.PerPage * (query.Page - 1)
	limit := query.PerPage

	err = db.Order("created_at DESC").
		Limit(limit).
		Offset(offSet).
		Find(&loans).Error
	return loans, count, err
}
