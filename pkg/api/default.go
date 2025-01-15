package api

import (
	"fmt"
	"github.com/bowoBp/LoanFlow/pkg/db"
	"github.com/gin-gonic/gin"
	"log"
)

func Default() *Api {
	server := gin.Default()
	_, err := db.Default()
	if err != nil {
		log.Println(err)
		panic(fmt.Sprintf("panic at db connection: %s", err.Error()))
	}
	fmt.Println("database connected: 3036")
	var routers = []Router{
		//user.NewRoute(sqlConn),
	}
	return &Api{
		server:  server,
		routers: routers,
	}
}
