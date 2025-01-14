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
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
)
