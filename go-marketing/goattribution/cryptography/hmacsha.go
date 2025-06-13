package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func HMACSHA256(buf, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func HMACSHA512(buf, key []byte) string {
	h := hmac.New(sha512.New, key)
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}
