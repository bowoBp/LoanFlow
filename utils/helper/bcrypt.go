package helper

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
}

type BcryptInterface interface {
	HasPass(pass string) string
	ComparePass(hashPass, pass []byte) bool
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
