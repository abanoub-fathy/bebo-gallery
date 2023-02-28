package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// Hasher is an object that can contain hmac inside it
// it can be used to make our code more easier
// while working with hashing
type Hasher struct {
	hmac hash.Hash
}

// NewHasher is used to create and return new Hasher
// and set the hmac secret key
func NewHasher(secretKey string) *Hasher {
	return &Hasher{
		hmac: hmac.New(sha256.New, []byte(secretKey)),
	}
}

// HashByHMAC is a method used to hash string and return the hashed string
// the hashing algorithm will use the secret key used when creating the hasher
func (h *Hasher) HashByHMAC(token string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(token))
	hashedByteSlice := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashedByteSlice)
}
