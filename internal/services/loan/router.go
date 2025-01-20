package loan

import (
	Repository "github.com/bowoBp/LoanFlow/internal/adapter/repository"
	"github.com/bowoBp/LoanFlow/internal/constant"
	"github.com/bowoBp/LoanFlow/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type (
	Router struct {
		auth middleware.AuthInterface
		rh   *RequestHandler
	}
)

func NewRoute(
	db *gorm.DB,
	auth middleware.AuthInterface,
) *Router {
	return &Router{
		auth: auth,
		rh: &RequestHandler{
			ctrl: &Controller{
				Uc: Usecase{
					LoanRepo:      Repository.NewLoanRepo(db),
					DbTransaction: NewLoanTransaction(db),
				},
			},
		},
	}
}

// Minimal handler to return 200 OK with "Not Implemented" message
func notImplementedHandler(c *gin.Context) {
	loanId := c.Param("loanId")
	log.Println("Loan ID:", loanId)
	c.JSON(http.StatusOK, gin.H{
		"message": "Not Implemented",
		"loanId":  loanId,
	})
}

func (r Router) Route(router *gin.RouterGroup) {
	loans := router.Group("loans")

	loans.GET(
		"/",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleAdmin,
			constant.RoleBorrower,
			constant.RoleStaff,
			constant.RoleInvestor,
		),
		r.rh.GetLoans,
	)
	loans.POST(
		"/",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleAdmin,
			constant.RoleBorrower,
		),
		r.rh.CreateLoan,
	)
	loans.POST(
		"/:loanId/approve",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleAdmin,
			constant.RoleStaff,
		),
		r.rh.ApproveLoan,
	)
	loans.POST(
		"/:loanId/invest",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleInvestor,
		),
		r.rh.StoreInvest,
	)
	loans.POST(
		"/:loanId/disburse",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleAdmin,
			constant.RoleStaff,
		),
		r.rh.DisburseLoan,
	)
	loans.GET(
		"/:loanId",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleAdmin,
			constant.RoleBorrower,
			constant.RoleStaff,
			constant.RoleInvestor,
		),
		r.rh.GetLoan,
	)
}
