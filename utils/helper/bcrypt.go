package helper

import (
	"crypto/hmac"
	"encoding/base32"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

type Bcrypt struct {
}

type BcryptInterface interface {
	HasPass(pass string) string
	ComparePass(hashPass, pass []byte) bool
	GenerateHashValue(
		secretKey string,
		uniqueID string,
		bitLen int,
	) (string, error)
}

func NewBcrypt() Bcrypt {
	return Bcrypt{}
}

func (r Bcrypt) HasPass(pass string) string {
	salt := 12
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), salt)
	return string(hashedPass)
}

func (r Bcrypt) ComparePass(hashPass, pass []byte) bool {
	hash, password := []byte(hashPass), []byte(pass)

	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

func (r Bcrypt) GenerateHashValue(
	secretKey string,
	uniqueID string,
	bitLen int,
) (string, error) {
	secretByte, err := base32.StdEncoding.DecodeString(secretKey)
	if err != nil {

		return "", err
	}
	hash := hmac.New(sha3.New224, secretByte)
	_, err = hash.Write([]byte(uniqueID))
	if err != nil {

		return "", err
	}
	hmacBytes := hash.Sum(nil)

	if bitLen > 1 {
		return hex.EncodeToString(hmacBytes[:bitLen]), nil
	}

	return hex.EncodeToString(hmacBytes), nil
}
