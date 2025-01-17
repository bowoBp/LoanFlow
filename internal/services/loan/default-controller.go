package loan

import (
	"context"
	"fmt"
	"github.com/bowoBp/LoanFlow/internal/constant"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"time"
)

type (
	Controller struct {
		Uc UsecaseInterface
	}

	ControllerInterface interface {
		CreateLoan(
			ctx context.Context,
			payload CreateLoanRequest,
		) (*dto.Response, error)
		GetLoan(
			ctx context.Context,
			loanID uint,
		) (*dto.Response, error)
		ApproveLoan(
			ctx context.Context,
			loanID, userID uint,
			payload ApproveLoanRequest,
		) (*dto.Response, error)
		DisburseLoan(
			ctx context.Context,
			loanID, userID uint,
			payload DisburseLoanRequest,
		) (*dto.Response, error)
		StoreInvest(
			ctx context.Context,
			loanID, userID uint,
			payload InvestLoanRequest,
		) (*dto.Response, error)
		GetLoans(
			ctx context.Context,
			query dto.GetListQuery,
		) (*dto.Response, error)
	}
)

func (ctrl Controller) CreateLoan(
	ctx context.Context,
	payload CreateLoanRequest,
) (*dto.Response, error) {
	start := time.Now()
	err := ctrl.Uc.CreateLoan(ctx, payload)
	if err != nil {
		return nil, err
	}
	return dto.NewSuccessResponse(
		payload,
		"success create loan",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil
}

func (ctrl Controller) ApproveLoan(
	ctx context.Context,
	loanID, userID uint,
	payload ApproveLoanRequest,
) (*dto.Response, error) {
	start := time.Now()
	err := ctrl.Uc.ApproveLoan(ctx, loanID, userID, payload)
	if err != nil {
		return nil, err
	}
	return dto.NewSuccessResponse(
		nil,
		"loan success approved",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil
}

func (ctrl Controller) StoreInvest(
	ctx context.Context,
	loanID, userID uint,
	payload InvestLoanRequest,
) (*dto.Response, error) {
	start := time.Now()
	err := ctrl.Uc.StoreInvest(ctx, loanID, userID, payload)
	if err != nil {
		return nil, err
	}
	return dto.NewSuccessResponse(
		nil,
		"success invested",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil
}

func (ctrl Controller) DisburseLoan(
	ctx context.Context,
	loanID, userID uint,
	payload DisburseLoanRequest,
) (*dto.Response, error) {
	start := time.Now()
	err := ctrl.Uc.DisburseLoan(ctx, loanID, userID, payload)
	if err != nil {
		return nil, err
	}
	return dto.NewSuccessResponse(
		nil,
		"loan success Disbursed",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil
}

func (ctrl Controller) GetLoan(
	ctx context.Context,
	loanID uint,
) (*dto.Response, error) {
	start := time.Now()
	res, err := ctrl.Uc.GetLoan(ctx, loanID)
	if err != nil {
		return nil, constant.LoanNotFound
	}
	result := &ResponseLoan{
		ID:              res.ID,
		BorrowerID:      res.BorrowerID,
		PrincipalAmount: res.PrincipalAmount,
		Rate:            res.Rate,
		Roi:             res.ROI,
		State:           res.State,
		AgreementLetter: res.AgreementLetterLink,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
	}

	return dto.NewSuccessResponse(
		result,
		"success get loan",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil
}

func (ctrl Controller) GetLoans(
	ctx context.Context,
	query dto.GetListQuery,
) (*dto.Response, error) {
	start := time.Now()
	if query.PerPage < 1 {
		query.PerPage = 10
	}
	if query.Page < 1 {
		query.Page = 1
	}
	loans, count, err := ctrl.Uc.GetLoans(ctx, query)
	if err != nil {
		return nil, err
	}
	result := make([]ListsLoanResponse, len(loans))
	for i, _ := range loans {
		result[i] = ListsLoanResponse{
			ID:              loans[i].ID,
			BorrowerID:      loans[i].BorrowerID,
			PrincipalAmount: loans[i].PrincipalAmount,
			Rate:            loans[i].Rate,
			Roi:             loans[i].ROI,
			State:           loans[i].State,
			CreatedAt:       loans[i].CreatedAt,
		}
	}

	paginatedResponse := dto.PaginationResponse{
		PerPage:     query.PerPage,
		Total:       uint(count),
		CurrentPage: query.Page,
	}
	paginatedResponse.Evaluate()
	listPaginated := ListLoanPaginate{
		Pagination: paginatedResponse,
		Data:       result,
	}

	return dto.NewSuccessResponse(
		listPaginated,
		"success get loan list",
		fmt.Sprint(time.Since(start).Milliseconds(), " ms."),
	), nil

}
