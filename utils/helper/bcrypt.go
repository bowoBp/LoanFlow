package helper

import "golang.org/x/crypto/bcrypt"

func HasPass(pass string) string {
	salt := 12
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), salt)
	return string(hashedPass)
}

func ComparePass(hashPass, pass []byte) bool {
	hash, password := []byte(hashPass), []byte(pass)

	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
