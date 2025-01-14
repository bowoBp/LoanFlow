package user

import (
	"github.com/bowoBp/LoanFlow/internal/adapter/Repository"
	"github.com/bowoBp/LoanFlow/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	Router struct {
		rq *RequestHandler
	}
)

func NewRoute(
	db *gorm.DB,
) *Router {
	return &Router{rq: &RequestHandler{
		ctrl: &Controller{
			Uc: UsaCase{
				userRepo: Repository.NewUserRepo(db),
			},
		},
	},
	}

}

func (r Router) Route(router *gin.RouterGroup) {
	employee := router.Group("/user")
	employee.POST(
		"/register",
		r.rq.Register,
	)
	employee.GET(
		"/all",
		middleware.Authentication(),
		r.rq.GetAll,
	)

	employee.GET(
		"/current",
		middleware.Authentication(),
		r.rq.GetCurrent,
	)

	employee.POST(
		"/login",
		r.rq.LoginCustomer,
	)

}
