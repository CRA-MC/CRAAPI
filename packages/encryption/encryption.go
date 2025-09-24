package encryption

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

type HashedPassword struct {
	Hashedpassword string
	Salt           string
	Iterations     int
}

func CreatePassword(plainpassword string) HashedPassword {
	var ret HashedPassword
	ret.Iterations = 8192
	salt := make([]byte, 32)
	rand.Read(salt)
	ret.Salt = base64.URLEncoding.EncodeToString(salt)
	ret.Hashedpassword = base64.URLEncoding.EncodeToString(pbkdf2.Key([]byte(plainpassword), []byte(ret.Salt), ret.Iterations, 256, sha256.New))
	return ret
}

func CheckPassowrd(plainpassword string, salt string, iterations int, checkpassword string) bool {

	if base64.URLEncoding.EncodeToString(pbkdf2.Key([]byte(plainpassword), []byte(salt), iterations, 256, sha256.New)) == checkpassword {
		return true
	} else {
		return false
	}
}
