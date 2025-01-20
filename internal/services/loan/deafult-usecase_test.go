package loan_test

import (
	"context"
	"errors"
	"github.com/bowoBp/LoanFlow/internal/adapter/mocks"
	"github.com/bowoBp/LoanFlow/internal/constant"
	domians "github.com/bowoBp/LoanFlow/internal/domain"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"github.com/bowoBp/LoanFlow/internal/services/loan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestUseCase_GetLoan(t *testing.T) {
	// Buat mock untuk LoanRepoInterface
	mockLoanRepo := new(mocks.LoanRepoInterface)

	// Mock data untuk test case
	mockLoan := &domians.Loan{
		ID:              1,
		BorrowerID:      42,
		PrincipalAmount: 1000000,
		Rate:            5.5,
		State:           "proposed",
	}

	type args struct {
		ctx    context.Context
		loanID uint
	}
	tests := []struct {
		name          string
		mockBehavior  func()
		args          args
		want          *domians.Loan
		wantErr       bool
		expectedError error
	}{
		{
			name: "Success - Loan Found",
			mockBehavior: func() {
				mockLoanRepo.On("GetLoanByID", mock.Anything, uint(1)).
					Return(mockLoan, nil).Once()
			},
			args: args{
				ctx:    context.Background(),
				loanID: 1,
			},
			want:    mockLoan,
			wantErr: false,
		},
		{
			name: "Error - Loan Not Found",
			mockBehavior: func() {
				mockLoanRepo.On("GetLoanByID", mock.Anything, uint(2)).
					Return(nil, errors.New("loan not found")).Once()
			},
			args: args{
				ctx:    context.Background(),
				loanID: 2,
			},
			want:          nil,
			wantErr:       true,
			expectedError: errors.New("loan not found"),
		},
		{
			name: "Error - Repository Failure",
			mockBehavior: func() {
				mockLoanRepo.On("GetLoanByID", mock.Anything, uint(3)).
					Return(nil, errors.New("repository failure")).Once()
			},
			args: args{
				ctx:    context.Background(),
				loanID: 3,
			},
			want:          nil,
			wantErr:       true,
			expectedError: errors.New("repository failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Atur perilaku mock sesuai test case
			tt.mockBehavior()

			// Buat instance Usecase
			uc := loan.Usecase{
				LoanRepo: mockLoanRepo,
			}

			// Panggil fungsi GetLoan
			got, err := uc.GetLoan(tt.args.ctx, tt.args.loanID)

			// Assert hasil
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)

			// Verifikasi bahwa mock dipanggil
			mockLoanRepo.AssertExpectations(t)
		})
	}
}

func TestUsecase_CreateLoan(t *testing.T) {

	// Mock dependencies
	mockTransaction := new(mocks.DefaultLoanTransactionInterface)

	// Mock data
	loanPayload := loan.CreateLoanRequest{
		ID:              1,
		PrincipalAmount: 1000000,
		Rate:            5.5,
	}
	mockLoan := &domians.Loan{
		ID:              1,
		BorrowerID:      1,
		PrincipalAmount: 1000000,
		Rate:            5.5,
		ROI:             55000,
		State:           "proposed",
	}

	type args struct {
		ctx     context.Context
		payload loan.CreateLoanRequest
	}

	tests := []struct {
		name         string
		mockBehavior func()
		args         args
		wantErr      bool
	}{
		{
			name: "Success - Loan Created",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock CreateLoan
				mockTransaction.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan *domians.Loan) bool {
					return loan != nil && loan.State == "proposed"
				})).Return(mockLoan, nil).Once()

				// Mock CreateLoanState
				mockTransaction.On("CreateLoanState", mock.Anything, mock.MatchedBy(func(state *domians.LoanStateHistory) bool {
					return state != nil && state.NewState == "proposed"
				})).Return(nil).Once()

				// Mock End
				mockTransaction.On("End", nil).Return(nil).Once()
			},
			args: args{
				ctx:     context.Background(),
				payload: loanPayload,
			},
			wantErr: false,
		},
		{
			name: "Error - Failed to Create Loan",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock CreateLoan with error
				mockTransaction.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan *domians.Loan) bool {
					return loan != nil && loan.State == "proposed"
				})).Return(nil, errors.New("failed to create loan")).Once()

				// Mock End with error
				mockTransaction.On("End", mock.MatchedBy(func(e error) bool {
					return e != nil && e.Error() == "failed to create loan"
				})).Return(nil).Once()
			},
			args: args{
				ctx:     context.Background(),
				payload: loanPayload,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock behavior
			tt.mockBehavior()

			// Create use case
			uc := loan.Usecase{
				DbTransaction: mockTransaction,
			}

			// Call CreateLoan
			err := uc.CreateLoan(tt.args.ctx, tt.args.payload)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all mocks
			mockTransaction.AssertExpectations(t)
		})
	}
}

func TestUsecase_ApproveLoan(t *testing.T) {
	// Mock dependencies
	mockTransaction := new(mocks.DefaultLoanTransactionInterface)
	mockLoanRepo := new(mocks.LoanRepoInterface)

	// Mock data
	loanID := uint(1)
	userID := uint(2)
	loanData := &domians.Loan{
		ID:    loanID,
		State: constant.Proposed,
	}

	payload := loan.ApproveLoanRequest{
		Proof: "link_to_proof",
	}

	type args struct {
		ctx     context.Context
		loanID  uint
		userID  uint
		payload loan.ApproveLoanRequest
	}

	tests := []struct {
		name         string
		mockBehavior func()
		args         args
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "Success - Loan Approved",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock GetLoanByID
				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).Return(loanData, nil).Once()

				// Mock UpdateLoan
				mockTransaction.On("UpdateLoan", mock.Anything, mock.MatchedBy(func(loan *domians.Loan) bool {
					return loan.ID == loanID
				}), mock.MatchedBy(func(data map[string]any) bool {
					return data["state"] == constant.Approved && data["agreement_letter_link"] == payload.Proof
				})).Return(nil).Once()

				// Mock ApproveDetail
				mockTransaction.On("ApproveDetail", mock.Anything, mock.MatchedBy(func(detail *domians.LoanApprovalDetail) bool {
					return detail.LoanID == loanID && detail.StaffID == userID && detail.PhotoProof == payload.Proof
				})).Return(nil).Once()

				// Mock CreateLoanState
				mockTransaction.On("CreateLoanState", mock.Anything, mock.MatchedBy(func(state *domians.LoanStateHistory) bool {
					return state.LoanID == loanID && state.PreviousState == constant.Proposed && state.NewState == constant.Approved
				})).Return(nil).Once()

				// Mock End
				mockTransaction.On("End", nil).Return(nil).Once()
			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr: false,
		},
		{
			name: "Error - Loan Not Found",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock GetLoanByID with error
				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(nil, errors.New("loan not found")).Once()

				// Mock End
				mockTransaction.On("End", mock.MatchedBy(func(e error) bool {
					return e != nil && e.Error() == constant.LoanNotFound.Error()
				})).Return(nil).Once()
			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr:     true,
			expectedErr: constant.LoanNotFound,
		},
		{
			name: "Error - Invalid Loan State",
			mockBehavior: func() {
				loanData.State = constant.Approved // Loan state is already approved

				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock GetLoanByID
				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(loanData, nil).Once()

				// Mock End
				mockTransaction.On("End", mock.MatchedBy(func(e error) bool {
					return e == nil || e.Error() == constant.ErrStateApprove.Error()
				})).Return(nil).Once()
			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr:     true,
			expectedErr: constant.ErrStateApprove,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock behavior
			tt.mockBehavior()

			// Create use case
			uc := loan.Usecase{
				LoanRepo:      mockLoanRepo,
				DbTransaction: mockTransaction,
			}

			// Call ApproveLoan
			err := uc.ApproveLoan(tt.args.ctx, tt.args.loanID, tt.args.userID, tt.args.payload)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all mocks
			mockTransaction.AssertExpectations(t)
			mockLoanRepo.AssertExpectations(t)
		})
	}
}

func TestUsecase_StoreInvest(t *testing.T) {
	// Mock dependencies
	mockTransaction := new(mocks.DefaultLoanTransactionInterface)
	mockLoanRepo := new(mocks.LoanRepoInterface)

	// Mock data
	loanID := uint(1)
	userID := uint(2)
	payload := loan.InvestLoanRequest{
		Amount: 150000,
	}

	type args struct {
		ctx     context.Context
		loanID  uint
		userID  uint
		payload loan.InvestLoanRequest
	}

	tests := []struct {
		name         string
		mockBehavior func()
		args         args
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "Success - Invested",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock GetLoanByID
				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(&domians.Loan{
						ID:              loanID,
						State:           constant.Approved,
						PrincipalAmount: 150000,
					}, nil).Once()

				// Mock UpdateLoan
				mockTransaction.On("UpdateLoan", mock.Anything, mock.MatchedBy(func(loan *domians.Loan) bool {
					return loan.ID == loanID
				}), mock.MatchedBy(func(data map[string]any) bool {
					return data["principal_amount"] == float64(0) && // Sesuaikan nilai aktual
						data["state"] == constant.Invested
				})).Return(nil).Once()

				// Mock InvestLoan
				mockTransaction.On("InvestLoan", mock.Anything, mock.MatchedBy(func(investor *domians.LoanInvestor) bool {
					return investor.LoanID == loanID &&
						investor.InvestorID == userID &&
						investor.AmountInvested == payload.Amount
				})).Return(nil).Once()

				// Mock CreateLoanState
				mockTransaction.On("CreateLoanState", mock.Anything, mock.MatchedBy(func(state *domians.LoanStateHistory) bool {
					return state.LoanID == loanID &&
						state.PreviousState == constant.Approved &&
						(state.NewState == "" || state.NewState == constant.Invested) &&
						state.ActionBy == userID
				})).Return(nil).Once()

				// Mock End
				mockTransaction.On("End", nil).
					Return(nil).Once()

			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr: false,
		},
		{
			name: "Error - Loan not found",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock GetLoanByID
				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(nil, errors.New("loan not found")).Once()

			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr:     true,
			expectedErr: constant.LoanNotFound,
		},
		{
			name: "Error - Invalid Loan State",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock GetLoanByID
				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(&domians.Loan{
						State: constant.Proposed,
					}, nil).Once()

			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr:     true,
			expectedErr: constant.ErrStateInvest,
		},
		{

			name: "Error - Invalid amount",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock GetLoanByID
				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(&domians.Loan{
						State:           constant.Approved,
						PrincipalAmount: 100000,
					}, nil).Once()

			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr:     true,
			expectedErr: constant.ErrInvestAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock behavior
			tt.mockBehavior()

			// Create use case
			uc := loan.Usecase{
				LoanRepo:      mockLoanRepo,
				DbTransaction: mockTransaction,
			}

			// Call StoreInvest
			err := uc.StoreInvest(tt.args.ctx, tt.args.loanID, tt.args.userID, tt.args.payload)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all mocks
			mockTransaction.AssertExpectations(t)
			mockLoanRepo.AssertExpectations(t)
		})
	}
}

func TestUsecase_DisburseLoan(t *testing.T) {
	// Mock dependencies
	mockTransaction := new(mocks.DefaultLoanTransactionInterface)
	mockLoanRepo := new(mocks.LoanRepoInterface)

	// Mock data
	loanID := uint(2)
	userID := uint(6)
	payload := loan.DisburseLoanRequest{
		AgreementLin: "agreement_link_file",
	}
	//loanData := &domians.Loan{
	//	ID:              loanID,
	//	State:           constant.Proposed, // State tidak sesuai ekspektasi
	//	PrincipalAmount: 1000000,
	//}

	type args struct {
		ctx     context.Context
		loanID  uint
		userID  uint
		payload loan.DisburseLoanRequest
	}
	tests := []struct {
		name         string
		mockBehavior func()
		args         args
		wantErr      bool
		expectedErr  error
	}{

		{
			name: "Error - loan not found",
			mockBehavior: func() {
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(nil, errors.New("loan not found")).Once()

			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr:     true,
			expectedErr: constant.LoanNotFound,
		},
		{
			name: "Error - state disbursed",
			mockBehavior: func() {

				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(&domians.Loan{
						ID:              loanID,
						PrincipalAmount: 1000000,
						State:           constant.Proposed,
					}, nil).Once()
			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr:     true,
			expectedErr: constant.ErrStateDisburse,
		},
		{
			name: "Success - Disburse",
			mockBehavior: func() {
				// Mock Begin
				mockTransaction.On("Begin").Return(mockTransaction, nil).Once()

				// Mock GetLoanByID
				mockLoanRepo.On("GetLoanByID", mock.Anything, loanID).
					Return(&domians.Loan{
						ID:    loanID,
						State: constant.Invested,
					}, nil).Once()

				// Mock UpdateLoan
				mockTransaction.On("UpdateLoan", mock.Anything, mock.MatchedBy(func(loan *domians.Loan) bool {
					return loan.ID == loanID
				}), mock.MatchedBy(func(data map[string]any) bool {
					return data["state"] == constant.Disbursed
				})).Return(nil).Once()

				// Mock InvestLoan
				mockTransaction.On("DisburseDetail", mock.Anything,
					mock.MatchedBy(func(disbursed *domians.LoanDisbursementDetail) bool {
						return disbursed.LoanID == loanID &&
							disbursed.StaffID == userID &&
							disbursed.SignedAgreementDoc == payload.AgreementLin
					})).Return(nil).Once()

				// Mock CreateLoanState
				mockTransaction.On("CreateLoanState", mock.Anything, mock.MatchedBy(func(state *domians.LoanStateHistory) bool {
					return state.LoanID == loanID &&
						state.PreviousState == constant.Invested &&
						state.NewState == constant.Disbursed &&
						state.ActionBy == userID &&
						state.Remarks == "Disbursed loan"
				})).Return(nil).Once()

				// Mock End
				mockTransaction.On("End", nil).
					Return(nil).Once()

			},
			args: args{
				ctx:     context.Background(),
				loanID:  loanID,
				userID:  userID,
				payload: payload,
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			uc := loan.Usecase{
				LoanRepo:      mockLoanRepo,
				DbTransaction: mockTransaction,
			}
			err := uc.DisburseLoan(tt.args.ctx, tt.args.loanID, tt.args.userID, tt.args.payload)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			mockTransaction.AssertExpectations(t)
			mockLoanRepo.AssertExpectations(t)
		})
	}
}

func TestUsecase_GetLoans(t *testing.T) {
	loanRepo := new(mocks.LoanRepoInterface)
	type args struct {
		ctx   context.Context
		query dto.GetListQuery
	}
	reqParam := dto.GetListQuery{
		PerPage: 10,
		Page:    1,
		Search:  "",
	}
	var LoanList = []domians.Loan{
		{
			ID:                  1,
			BorrowerID:          101,
			PrincipalAmount:     500000,
			Rate:                7.5,
			ROI:                 37500,
			State:               "proposed",
			AgreementLetterLink: "",
			CreatedAt:           time.Now(),
		},
		{
			ID:                  2,
			BorrowerID:          102,
			PrincipalAmount:     1000000,
			Rate:                5.5,
			ROI:                 55000,
			State:               "approved",
			AgreementLetterLink: "http://example.com/agreement_2.pdf",
			CreatedAt:           time.Now(),
		},
	}
	tests := []struct {
		name         string
		mockBehavior func()
		args         args
		want         []domians.Loan
		want1        int64
		wantErr      bool
		expectedErr  error
	}{
		{
			name: "success - get list loan",
			mockBehavior: func() {
				loanRepo.On("GetLoans", mock.Anything, reqParam).
					Return(LoanList, int64(50), nil).Once()
			},
			args: args{
				ctx:   context.Background(),
				query: reqParam,
			},
			want:        LoanList,
			want1:       50,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			uc := loan.Usecase{
				LoanRepo: loanRepo,
			}
			_, _, err := uc.GetLoans(tt.args.ctx, tt.args.query)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			loanRepo.AssertExpectations(t)
		})
	}
}
