package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

const secretKey = "access-login"

type JwtPayload struct {
	ID        string
	UserName  string
	CreatedAt time.Time
}

func GenerateToken(id, userName string, createdAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"id":        id,
		"userName":  userName,
		"createdAt": createdAt,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(c *gin.Context) (string, error) {
	errResponse := errors.New("token invalid")
	headerToken := c.Request.Header.Get("Authorization")

	// Memeriksa apakah token memiliki prefix "Bearer "
	if !strings.HasPrefix(headerToken, "Bearer ") {
		return "", errResponse
	}

	// Memisahkan "Bearer" dari token yang sebenarnya
	parts := strings.Split(headerToken, " ")
	if len(parts) != 2 {
		return "", errResponse // Mengembalikan error jika format token tidak benar
	}
	stringToken := parts[1]

	// Parse token menggunakan secret key
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", errResponse // Mengembalikan error jika token tidak dapat diparse atau method salah
	}

	// Memeriksa kevalidan token
	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return "", errResponse
	} else {
		// Jika token valid, Anda bisa mengakses claims-nya di sini jika diperlukan
		fmt.Println(claims) // Opsional: Menampilkan claims
	}

	return stringToken, nil
}

func ExtractPayloadFromToken(requestToken string) (res JwtPayload, err error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return res, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return res, errors.New("invalid token")
	}

	cb, _ := json.Marshal(claims)
	json.Unmarshal(cb, &res)

	return res, nil
}
