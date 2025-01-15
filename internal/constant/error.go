package constant

import "errors"

var (
	ErrRegister     = errors.New("password or username is required")
	ErrUserNotFound = errors.New("user not found")
	ErrPriceIsMinus = errors.New("not allowed negative price")
	ErrProductName  = errors.New("product name can't empty")
	DuplicateEmail  = errors.New("duplicate email")
)
