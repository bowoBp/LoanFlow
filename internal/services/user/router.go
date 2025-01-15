package user

import (
	"github.com/bowoBp/LoanFlow/internal/adapter/Repository"
	"github.com/bowoBp/LoanFlow/pkg/environment"
	"github.com/bowoBp/LoanFlow/pkg/middleware"
	"github.com/bowoBp/LoanFlow/utils/helper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	Router struct {
		rq   *RequestHandler
		auth middleware.AuthInterface
	}
)

func NewRoute(
	db *gorm.DB,
	jwt helper.JwtInterface,
	bcrypt helper.BcryptInterface,
	env environment.Environment,
	auth middleware.AuthInterface,
) *Router {
	return &Router{
		rq: &RequestHandler{
			ctrl: &Controller{
				Uc: UsaCase{
					userRepo: Repository.NewUserRepo(db),
					jwt:      jwt,
					bcrypt:   bcrypt,
					env:      env,
				},
			},
		},
		auth: auth,
	}

}

func (r Router) Route(router *gin.RouterGroup) {
	employee := router.Group("/user")
	auth := router.Group("/auth")
	employee.POST(
		"/register",
		r.rq.Register,
	)
	employee.GET(
		"/all",
		r.auth.Authentication(),
		r.rq.GetAll,
	)

	employee.GET(
		"/current",
		r.auth.Authentication(),
		r.rq.GetCurrent,
	)

	auth.POST(
		"/login",
		r.rq.Login,
	)

}
