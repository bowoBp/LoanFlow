package constant

import "errors"

var (
	ErrRegister      = errors.New("password or username is required")
	ErrUserNotFound  = errors.New("user not found")
	ErrPriceIsMinus  = errors.New("not allowed negative price")
	ErrProductName   = errors.New("product name can't empty")
	DuplicateEmail   = errors.New("duplicate email")
	LoanNotFound     = errors.New("loan not found")
	ErrStateApprove  = errors.New("only loans in 'proposed' state can be approved")
	ErrStateDisburse = errors.New("only loans in 'Invested' state can be disburse")
	ErrStateInvest   = errors.New("loan is not in 'approved' state'")
	ErrInvestAmount  = errors.New("investment amount exceeds the principal amount")
)
