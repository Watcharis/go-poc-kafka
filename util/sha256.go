package util

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func SHA256Encode(in string) string {
	data := []byte(in)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:])
}

func SHA512Encode(in string) string {
	data := []byte(in)
	hash := sha512.Sum512(data)
	return fmt.Sprintf("%x", hash[:])
}
