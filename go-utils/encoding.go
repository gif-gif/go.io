package goutils

import (
	"bytes"
	golog "github.com/gif-gif/go.io/go-log"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func GBK2UTF8(s string) string {
	r := bytes.NewBufferString(s)
	rr := transform.NewReader(r, simplifiedchinese.GBK.NewDecoder())
	buf, err := ioutil.ReadAll(rr)
	if err != nil {
		golog.WithField("str", s).Error(err.Error())
		return ""
	}
	return string(bytes.TrimSpace(buf))
}

func UTF82GBK(s string) string {
	r := bytes.NewBufferString(s)
	rr := transform.NewReader(r, simplifiedchinese.GBK.NewEncoder())
	buf, err := ioutil.ReadAll(rr)
	if err != nil {
		golog.WithField("str", s).Error(err.Error())
		return ""
	}
	return string(bytes.TrimSpace(buf))
}
