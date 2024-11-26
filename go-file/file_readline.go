package gofile

import (
	"bufio"
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
	"io"
	"os"
)

func ReadByLine(filename string, cb func(b []byte, end bool) error) error {
	if !Exists(filename) {
		return errors.New("文件不存在")
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0755)
	if err != nil {
		golog.Error(err)
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, err := r.ReadBytes('\n')

		if err != nil {
			if io.EOF == err {
				return cb(b, true)
			}

			golog.Error(err)
			return err
		}

		if err := cb(b, false); err != nil {
			return err
		}
	}
}
