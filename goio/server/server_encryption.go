package goserver

import (
	"encoding/hex"

	"github.com/gif-gif/go.io/go-utils/gocrypto"

	"strings"
)

type Encryption struct {
	Key    string
	Secret string
}

func (enc *Encryption) Encode(b []byte) (str string, err error) {
	if l := len(b); l == 0 {
		return
	}
	var bts []byte
	bts, err = gocrypto.AESCBCEncrypt(b, []byte(enc.Key), []byte(enc.Secret))
	if err != nil {
		return
	}
	str = hex.EncodeToString(bts)
	return
}

func (enc *Encryption) Decode(str string) (b []byte, err error) {
	str = strings.ReplaceAll(str, "\"", "")
	if l := len(str); l == 0 {
		return
	}
	var bts []byte
	bts, err = hex.DecodeString(str)
	if err != nil {
		return
	}
	b, err = gocrypto.AESCBCDecrypt(bts, []byte(enc.Key), []byte(enc.Secret))
	return
}
