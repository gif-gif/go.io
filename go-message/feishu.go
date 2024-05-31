package gomessage

import (
	"errors"
	gohttp "github.com/gif-gif/go.io/go-http"
	goutils "github.com/gif-gif/go.io/go-utils"
	"runtime"
	"sync"
)

var (
	__fieShuCH   chan struct{}
	__feiShuOnce sync.Once
)

func FeiShu(hookUrl string, text string) error {
	__feiShuOnce.Do(func() {
		__fieShuCH = make(chan struct{}, runtime.NumCPU()*2)
	})

	__fieShuCH <- struct{}{}
	defer func() { <-__fieShuCH }()

	content := goutils.NewParams().Set("text", text)

	params := goutils.NewParams().
		Set("msg_type", "text").
		Set("content", content.Data())

	buf, err := gohttp.PostJson(hookUrl, params.JSON())
	if err != nil {
		return err
	}

	rst, _ := goutils.Byte(buf).Params()
	if msg := rst.Get("StatusMessage").String(); msg != "success" {
		return errors.New(msg)
	}

	return nil
}
