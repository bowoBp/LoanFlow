package loan

import (
	"context"
	Repository "github.com/bowoBp/LoanFlow/internal/adapter/repository"
	"github.com/bowoBp/LoanFlow/internal/constant"
	domians "github.com/bowoBp/LoanFlow/internal/domain"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"time"
)

type (
	UseCase struct {
		loanRepo      Repository.LoanRepoInterface
		dbTransaction Repository.TransactionUnit[DefaultLoanTransactionInterface]
	}

	UsecaseInterface interface {
		CreateLoan(
			ctx context.Context,
			payload CreateLoanRequest,
		) error

		GetLoan(
			ctx context.Context,
			loanID uint,
		) (*domians.Loan, error)
		ApproveLoan(
			ctx context.Context,
			loanID, userID uint,
			payload ApproveLoanRequest,
		) error
		DisburseLoan(
			ctx context.Context,
			loanID, userID uint,
			payload DisburseLoanRequest,
		) error
		StoreInvest(
			ctx context.Context,
			loanID, userID uint,
			payload InvestLoanRequest,
		) error
		GetLoans(
			ctx context.Context,
			query dto.GetListQuery,
		) ([]domians.Loan, int64, error)
	}
)

func (uc UseCase) CreateLoan(
	ctx context.Context,
	payload CreateLoanRequest,
) error {
	dbTrx, err := uc.dbTransaction.Begin()
	defer func(tx DefaultLoanTransactionInterface, err *error) {
		// recover panic
		if r := recover(); r != nil {
			// TODO: catch error and pass to log/sentry soon
		}
		// end transaction (rollback or commit)
		errTrx := tx.End(*err)
		if errTrx != nil {
			// TODO: catch error and pass to log/sentry soon
		}
	}(dbTrx, &err)
	loan, err := dbTrx.CreateLoan(
		ctx,
		&domians.Loan{
			BorrowerID:      uint(payload.ID),
			PrincipalAmount: payload.PrincipalAmount,
			Rate:            payload.Rate,
			ROI:             payload.PrincipalAmount * (payload.Rate / 100),
			State:           constant.Proposed,
		},
	)
	if err != nil {
		return err
	}

	err = dbTrx.CreateLoanState(
		ctx,
		&domians.LoanStateHistory{
			LoanID:        loan.ID,
			PreviousState: "",
			NewState:      constant.Proposed,
			ActionBy:      payload.ID,
			ActionAt:      time.Now(),
			Remarks:       "proposed loan",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (uc UseCase) GetLoan(
	ctx context.Context,
	loanID uint,
) (*domians.Loan, error) {
	return uc.loanRepo.GetLoanByID(ctx, loanID)
}

func (uc UseCase) ApproveLoan(
	ctx context.Context,
	loanID, userID uint,
	payload ApproveLoanRequest,
) error {
	dbTrx, err := uc.dbTransaction.Begin()
	defer func(tx DefaultLoanTransactionInterface, err *error) {
		// recover panic
		if r := recover(); r != nil {
			// TODO: catch error and pass to log/sentry soon
		}
		// end transaction (rollback or commit)
		errTrx := tx.End(*err)
		if errTrx != nil {
			// TODO: catch error and pass to log/sentry soon
		}
	}(dbTrx, &err)

	loan, err := uc.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return constant.LoanNotFound
	}
	if loan.State != constant.Proposed {
		return constant.ErrStateApprove
	}
	err = dbTrx.UpdateLoan(
		ctx,
		&domians.Loan{
			ID: loanID,
		},
		map[string]any{
			"state":                 constant.Approved,
			"agreement_letter_link": payload.Proof,
			"updated_at":            time.Now(),
		},
	)
	if err != nil {

		return err
	}
	err = dbTrx.ApproveDetail(
		ctx,
		&domians.LoanApprovalDetail{
			LoanID:       loanID,
			StaffID:      userID,
			PhotoProof:   payload.Proof,
			ApprovedDate: time.Now(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	if err != nil {
		return err
	}
	err = dbTrx.CreateLoanState(
		ctx,
		&domians.LoanStateHistory{
			LoanID:        loanID,
			PreviousState: loan.State,
			NewState:      constant.Approved,
			ActionBy:      userID,
			ActionAt:      time.Now(),
			Remarks:       "Approved loan",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (uc UseCase) StoreInvest(
	ctx context.Context,
	loanID, userID uint,
	payload InvestLoanRequest,
) error {
	dbTrx, err := uc.dbTransaction.Begin()
	nextState := ""
	defer func(tx DefaultLoanTransactionInterface, err *error) {
		// recover panic
		if r := recover(); r != nil {
			// TODO: catch error and pass to log/sentry soon
		}
		// end transaction (rollback or commit)
		errTrx := tx.End(*err)
		if errTrx != nil {
			// TODO: catch error and pass to log/sentry soon
		}
	}(dbTrx, &err)

	loan, err := uc.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return constant.LoanNotFound
	}
	if loan.State != constant.Approved {
		return constant.ErrStateInvest
	}
	if loan.PrincipalAmount < payload.Amount {
		return constant.ErrInvestAmount
	}
	loan.PrincipalAmount -= payload.Amount
	if loan.PrincipalAmount == 0 {
		loan.State = constant.Invested
		nextState = constant.Invested
	}
	err = dbTrx.UpdateLoan(
		ctx,
		&domians.Loan{
			ID: loanID,
		},
		map[string]any{
			"state":            loan.State,
			"principal_amount": loan.PrincipalAmount,
			"updated_at":       time.Now(),
		},
	)
	if err != nil {

		return err
	}
	err = dbTrx.InvestLoan(
		ctx,
		&domians.LoanInvestor{
			LoanID:         loanID,
			InvestorID:     userID,
			AmountInvested: payload.Amount,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	if err != nil {
		return err
	}

	err = dbTrx.CreateLoanState(
		ctx,
		&domians.LoanStateHistory{
			LoanID:        loanID,
			PreviousState: constant.Approved,
			NewState:      nextState,
			ActionBy:      userID,
			ActionAt:      time.Now(),
			Remarks:       "user invested",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (uc UseCase) DisburseLoan(
	ctx context.Context,
	loanID, userID uint,
	payload DisburseLoanRequest,
) error {
	dbTrx, err := uc.dbTransaction.Begin()
	defer func(tx DefaultLoanTransactionInterface, err *error) {
		// recover panic
		if r := recover(); r != nil {
			// TODO: catch error and pass to log/sentry soon
		}
		// end transaction (rollback or commit)
		errTrx := tx.End(*err)
		if errTrx != nil {
			// TODO: catch error and pass to log/sentry soon
		}
	}(dbTrx, &err)
	loan, err := uc.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return constant.LoanNotFound
	}
	if loan.State != constant.Invested {
		return constant.ErrStateDisburse
	}
	err = dbTrx.UpdateLoan(
		ctx,
		&domians.Loan{
			ID: loanID,
		},
		map[string]any{
			"state":      constant.Disbursed,
			"updated_at": time.Now(),
		},
	)
	if err != nil {

		return err
	}
	err = dbTrx.DisburseDetail(
		ctx, &domians.LoanDisbursementDetail{
			LoanID:             loanID,
			StaffID:            userID,
			SignedAgreementDoc: payload.AgreementLin,
			DisbursedDate:      time.Now(),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		})
	if err != nil {
		return err
	}
	err = dbTrx.CreateLoanState(
		ctx,
		&domians.LoanStateHistory{
			LoanID:        loanID,
			PreviousState: loan.State,
			NewState:      constant.Disbursed,
			ActionBy:      userID,
			ActionAt:      time.Now(),
			Remarks:       "Disbursed loan",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (uc UseCase) GetLoans(
	ctx context.Context,
	query dto.GetListQuery,
) ([]domians.Loan, int64, error) {
	return uc.loanRepo.GetLoans(ctx, query)
}
