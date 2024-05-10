package goutils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestId2Code(t *testing.T) {
	for i := 1; i < 1000; i++ {
		code := Id2Code(int64(i))
		id, _ := Code2Id(code)
		fmt.Println(id, code)
	}
}

func TestKey(t *testing.T) {
	key := []rune("123567890ABCDEFGHJKLMNPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())

	var data []rune

	for {
		l := len(key)
		if l == 0 {
			break
		}
		if l == 1 {
			data = append(data, key[0])
			break
		}

		n := rand.Intn(l - 1)
		data = append(data, key[n])

		key = append(key[:n], key[n+1:]...)
	}

	fmt.Println(string(data))
}
