package goio

import (
	"log"
	"testing"
)

func TestEncryption_Encode(t *testing.T) {
	enc := &Encryption{
		Key:    "4c98542af9fd65fc",
		Secret: "12fa1e087ab15eba558e12ea64d0f3f8",
	}

	str, err := enc.Encode([]byte("123adbasdf"))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(str)
}

func TestEncryption_Decode(t *testing.T) {
	enc := &Encryption{
		Key:    "4c98542af9fd65fc",
		Secret: "12fa1e087ab15eba558e12ea64d0f3f8",
	}

	b, err := enc.Decode("a12241192af531a8361b9d195bb9a7863ae5a9507bbb3ad589b6f18293e18509df3bd1174e8eea1afc6ac73d1b70ff26")
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(string(b))
}
