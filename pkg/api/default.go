package api

import (
	"fmt"
	"github.com/bowoBp/LoanFlow/internal/services/loan"
	"github.com/bowoBp/LoanFlow/internal/services/user"
	"github.com/bowoBp/LoanFlow/pkg/db"
	"github.com/bowoBp/LoanFlow/pkg/environment"
	"github.com/bowoBp/LoanFlow/pkg/middleware"
	"github.com/bowoBp/LoanFlow/utils/helper"
	"github.com/gin-gonic/gin"
	"log"
)

func Default() *Api {
	server := gin.Default()
	sqlConn, err := db.Default()
	jwt := helper.NewJwt()
	bcrypt := helper.NewBcrypt()
	env := environment.NewEnvironment()
	auth := middleware.NewAuth()
	if err != nil {
		log.Println(err)
		panic(fmt.Sprintf("panic at db connection: %s", err.Error()))
	}
	fmt.Println("database connected: 3036")
	var routers = []Router{
		user.NewRoute(sqlConn, jwt, bcrypt, env, auth),
		loan.NewRoute(auth),
	}
	return &Api{
		server:  server,
		routers: routers,
	}
}
