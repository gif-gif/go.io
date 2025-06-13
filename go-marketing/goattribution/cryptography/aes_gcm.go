package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AesGcmDecrypt(cipherText, key, nonce []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(nonce) < nonceSize {
		return nil, errors.New("cipherText too short")
	}

	// nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	return gcm.Open(nil, nonce, cipherText, nil)
}
