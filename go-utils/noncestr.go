package goutils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func NonceStr() string {
	bf := make([]byte, 8)
	io.ReadFull(rand.Reader, bf)
	return hex.EncodeToString(bf)
}

func NonceStr8() string {
	bf := make([]byte, 4)
	io.ReadFull(rand.Reader, bf)
	return hex.EncodeToString(bf)
}
