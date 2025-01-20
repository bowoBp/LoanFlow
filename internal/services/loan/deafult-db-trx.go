package loan

import (
	"context"
	Repository "github.com/bowoBp/LoanFlow/internal/adapter/repository"
	domians "github.com/bowoBp/LoanFlow/internal/domain"
	"gorm.io/gorm"
)

type (
	DefaultLoanTransaction struct {
		db       *gorm.DB
		loanRepo Repository.LoanRepoInterface
	}

	DefaultLoanTransactionInterface interface {
		Begin() (DefaultLoanTransactionInterface, error)
		End(err error) error
		CreateLoan(
			ctx context.Context,
			loan *domians.Loan,
		) (*domians.Loan, error)
		ApproveDetail(
			ctx context.Context,
			detail *domians.LoanApprovalDetail,
		) error
		DisburseDetail(
			ctx context.Context,
			disbursed *domians.LoanDisbursementDetail,
		) error
		CreateLoanState(
			ctx context.Context,
			state *domians.LoanStateHistory,
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
	}
)

func NewLoanTransaction(db *gorm.DB) DefaultLoanTransaction {
	return DefaultLoanTransaction{
		db: db,
	}
}

func (repo DefaultLoanTransaction) UpdateLoan(
	ctx context.Context,
	loan *domians.Loan,
	updateData map[string]any,
) error {
	return repo.loanRepo.UpdateLoan(ctx, loan, updateData)
}

func (repo DefaultLoanTransaction) ApproveDetail(
	ctx context.Context,
	detail *domians.LoanApprovalDetail,
) error {
	return repo.loanRepo.ApproveDetail(ctx, detail)
}

func (repo DefaultLoanTransaction) CreateLoanState(
	ctx context.Context,
	state *domians.LoanStateHistory,
) error {
	return repo.loanRepo.CreateLoanState(ctx, state)
}

func (repo DefaultLoanTransaction) CreateLoan(
	ctx context.Context,
	loan *domians.Loan,
) (*domians.Loan, error) {
	return repo.loanRepo.CreateLoan(ctx, loan)
}

func (repo DefaultLoanTransaction) DisburseDetail(
	ctx context.Context,
	disbursed *domians.LoanDisbursementDetail,
) error {
	return repo.loanRepo.DisburseDetail(ctx, disbursed)
}

func (repo DefaultLoanTransaction) InvestLoan(
	ctx context.Context,
	investLoan *domians.LoanInvestor,
) error {
	return repo.loanRepo.InvestLoan(ctx, investLoan)
}

func (repo DefaultLoanTransaction) Begin() (DefaultLoanTransactionInterface, error) {
	evoTrx := repo.db.Begin()
	if err := evoTrx.Error; err != nil {
		return DefaultLoanTransaction{}, err
	}
	newLoanTrx := &DefaultLoanTransaction{
		db:       evoTrx,
		loanRepo: Repository.NewLoanRepo(evoTrx),
	}
	return newLoanTrx, nil
}

func (repo DefaultLoanTransaction) End(err error) error {
	if err != nil {
		errTrx := repo.db.Rollback().Error
		if errTrx != nil {
			return nil
		}
	}
	errTrx := repo.db.Commit().Error
	if errTrx != nil {
		return errTrx
	}
	return nil
}
