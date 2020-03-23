package library

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt ... Tạo ra password mã hoá có kèm theo muối
func HashAndSalt(pwd string) string {
	byteString := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(byteString, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
