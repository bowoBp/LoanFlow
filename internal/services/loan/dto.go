package loan

import (
	"github.com/bowoBp/LoanFlow/internal/dto"
	"time"
)

type (
	CreateLoanRequest struct {
		ID              uint    `json:"id"`
		PrincipalAmount float64 `json:"principalAmount"`
		Rate            float64 `json:"rate"`
	}

	ResponseLoan struct {
		ID              uint      `json:"id"`
		BorrowerID      uint      `json:"borrowerId"`
		PrincipalAmount float64   `json:"principalAmount"`
		Rate            float64   `json:"rate"`
		Roi             float64   `json:"roi"`
		State           string    `json:"state"`
		AgreementLetter string    `json:"agreementLetter"`
		CreatedAt       time.Time `json:"createdAt"`
		UpdatedAt       time.Time `json:"updatedAt"`
	}

	ListsLoanResponse struct {
		ID              uint      `json:"id"`
		BorrowerID      uint      `json:"borrowerId"`
		PrincipalAmount float64   `json:"principalAmount"`
		Rate            float64   `json:"rate"`
		Roi             float64   `json:"roi"`
		State           string    `json:"state"`
		CreatedAt       time.Time `json:"createdAt"`
	}
	ListLoanPaginate struct {
		Pagination dto.PaginationResponse
		Data       []ListsLoanResponse `json:"data"`
	}

	ApproveLoanRequest struct {
		Proof string `json:"proof"`
	}

	DisburseLoanRequest struct {
		AgreementLin string `json:"agreementLink"`
	}

	InvestLoanRequest struct {
		Amount float64 `json:"amount"`
	}
)
