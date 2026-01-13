package crypt

import (
	"fmt"

	"golang.org/x/crypto/scrypt"
)

var salt = "xhc"

func PasswordEncrypt(data string) string {
	dk, _ := scrypt.Key([]byte(data), []byte(salt), 1<<15, 8, 1, 32)
	return fmt.Sprintf("%x", string(dk))
}
