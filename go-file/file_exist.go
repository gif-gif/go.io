package gofile

import (
	"os"
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}
