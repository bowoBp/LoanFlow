package user

import (
	"github.com/bowoBp/LoanFlow/internal/dto"
	"time"
)

type (
	RegisterUser struct {
		ID        string    `json:"id"`
		UserName  string    `validate:"required" json:"userName"`
		Password  string    `validate:"required" json:"password"`
		Email     string    `validate:"required" json:"email"`
		Phone     string    `validate:"required" json:"phone"`
		Role      string    `validate:"required" json:"role"`
		CreatedAt time.Time `json:"createdAt"`
	}

	UseCaseRegisterResult struct {
		User RegisterUser `json:"user"`
	}

	SuccessLoginUser struct {
		Response    dto.ResponseMeta
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	}

	Users struct {
		ID        string    `json:"id"`
		UserName  string    `json:"UserName"`
		CreatedAt time.Time `json:"createdAt"`
	}

	LoginParam struct {
		Email    string `validate:"required" json:"email"`
		Password string `validate:"required" json:"password"`
	}
)
