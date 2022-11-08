package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenSize = 32

func RandBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func RandString(n int) (string, error) {
	randomBytes, err := RandBytes(n)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}


// GenerateRememberToken is used to generate string remeber token with predefined size
// this token is not hashed yet
func GenerateRememberToken() (string, error){
	return RandString(RememberTokenSize)
}