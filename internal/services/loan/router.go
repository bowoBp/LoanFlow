package loan

import (
	"github.com/bowoBp/LoanFlow/internal/constant"
	"github.com/bowoBp/LoanFlow/pkg/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type (
	Router struct {
		auth middleware.AuthInterface
	}
)

func NewRoute(
	auth middleware.AuthInterface,
) *Router {

	return &Router{
		auth: auth,
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
		notImplementedHandler,
	)
	loans.POST(
		"/",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleAdmin,
			constant.RoleBorrower,
		),
		notImplementedHandler,
	)
	loans.POST(
		"/:loanId/approve",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleAdmin,
			constant.RoleStaff,
		),
		notImplementedHandler,
	)
	loans.POST(
		"/:loanId/invest",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleInvestor,
		),
		notImplementedHandler,
	)
	loans.POST(
		"/:loanId/disburse",
		r.auth.Authentication(),
		r.auth.Authorize(
			constant.RoleAdmin,
			constant.RoleStaff,
		),
		notImplementedHandler,
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
		notImplementedHandler,
	)
}
