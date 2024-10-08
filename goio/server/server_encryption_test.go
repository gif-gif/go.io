package goserver

import (
	"log"
	"testing"
)

func TestEncryption_Encode(t *testing.T) {
	enc := &Encryption{
		Key:    "18586555d498e2c3d0e0baaf0bdde3ce",
		Secret: "7732c2d83bad98173790cd62483cde791df5ce6e62247dc563a15367e4750340",
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
		Key:    "f455c40c0303189b7a1fe769f7689a52",
		Secret: "aa4ef454a2c8d157fa5ea43df5357806",
	}

	b, err := enc.Decode("269e4be01b92d9f4c0a6477a4e5a3b49")
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(string(b))
}
