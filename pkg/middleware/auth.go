package middleware

import (
	"fmt"
	"github.com/bowoBp/LoanFlow/internal/dto"
	"github.com/bowoBp/LoanFlow/utils/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"time"
)

type Auth struct {
	jwt helper.JwtInterface
}

type AuthInterface interface {
	Authentication() gin.HandlerFunc
	Authorize(roles ...string) gin.HandlerFunc
}

func NewAuth() AuthInterface {
	return &Auth{
		jwt: helper.NewJwt(),
	}
}

func (receiver Auth) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		stringToken, err := receiver.jwt.VerifyToken(c)
		if err != nil {
			response := dto.DefaultErrorResponseWithMessage(err.Error())
			response.ResponseTime = fmt.Sprint(time.Since(start).Milliseconds(), " ms.")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		jwtPayload, _ := receiver.jwt.ExtractPayloadFromToken(stringToken)

		c.Set("id", jwtPayload.ID)
		c.Set("userName", jwtPayload.UserName)
		c.Set("createdAt", jwtPayload.CreatedAt)
		c.Set("userRole", jwtPayload.UserRole)
		c.Next()
	}
}

func (receiver Auth) Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Ambil userRole dari context
		valUserRole, ok := c.Get("userRole")
		if !ok {
			response := dto.DefaultBadRequestResponse()
			response.Message = "invalid token"
			response.ResponseTime = fmt.Sprintf("%d ms.", time.Since(start).Milliseconds())
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// Pastikan userRole adalah string
		userRole, ok := valUserRole.(string)
		if !ok {
			response := dto.DefaultBadRequestResponse()
			response.Message = "invalid token"
			response.ResponseTime = fmt.Sprintf("%d ms.", time.Since(start).Milliseconds())
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// Periksa apakah role yang digunakan user ada di dalam list roles yang diizinkan
		if slices.Contains(roles, userRole) {
			c.Next()
			return
		}

		// Jika tidak ada di dalam roles yang diizinkan
		response := dto.DefaultBadRequestResponse()
		response.Message = "Kamu tidak punya akses ke halaman ini"
		response.ResponseTime = fmt.Sprintf("%d ms.", time.Since(start).Milliseconds())
		c.JSON(http.StatusUnauthorized, response)
		c.Abort()
	}
}
