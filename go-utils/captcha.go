package goutils

import (
	"bytes"
	"fmt"
	"github.com/dchest/captcha"
)

// 获取图片验证码
func CaptchaGet(width, height int) map[string]string {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}

	var buf bytes.Buffer

	id := captcha.NewLen(4)
	captcha.WriteImage(&buf, id, width, height)

	b64 := fmt.Sprintf("data:image/png;base64,%s", Base64Encode(buf.Bytes()))

	return map[string]string{
		"id":          id,
		"base64image": b64,
	}
}

// 验证图片验证码
func CaptchaVerify(id, code string) bool {
	return captcha.VerifyString(id, code)
}
