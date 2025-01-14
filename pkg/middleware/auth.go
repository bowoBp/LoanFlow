package middleware

import (
	"fmt"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"github.com/bowoBp/LoanFlow/utils/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		stringToken, err := helper.VerifyToken(c)
		if err != nil {
			response := dto.DefaultErrorResponseWithMessage(err.Error())
			response.ResponseTime = fmt.Sprint(time.Since(start).Milliseconds(), " ms.")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		jwtPayload, _ := helper.ExtractPayloadFromToken(stringToken)

		c.Set("id", jwtPayload.ID)
		c.Set("userName", jwtPayload.UserName)
		c.Set("createdAt", jwtPayload.CreatedAt)
		c.Next()
	}
}
