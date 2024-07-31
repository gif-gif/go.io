package goutils

import (
	gocaptcha "github.com/gif-gif/go.io/go-captcha"
)

//TODO: https://github.com/mojocn/base64Captcha
//https://github.com/wenlng/go-captcha (高级)

var defaultCaptcha = gocaptcha.NewDefault()

func InitCaptcha(captcha *gocaptcha.GoCaptcha) {
	defaultCaptcha = captcha
}

// 获取图片验证码-数字
func CaptchaGet(width, height int) map[string]string {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}

	data, err := defaultCaptcha.DigitCaptcha(width, height, 4)
	if err != nil {
		return map[string]string{
			"id":          "",
			"base64image": "",
		}
	}

	return map[string]string{
		"id":          data.CaptchaId,
		"base64image": data.Data,
	}
}

func CaptchaChineseGet(width, height int, source string) map[string]string {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}

	data, err := defaultCaptcha.ChineseCaptcha(width, height, 4, source)
	if err != nil {
		return map[string]string{
			"id":          "",
			"base64image": "",
		}
	}

	return map[string]string{
		"id":          data.CaptchaId,
		"base64image": data.Data,
	}
}

func CaptchaMathGet(width, height int) map[string]string {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}

	data, err := defaultCaptcha.MathCaptcha(width, height)
	if err != nil {
		return map[string]string{
			"id":          "",
			"base64image": "",
		}
	}

	return map[string]string{
		"id":          data.CaptchaId,
		"base64image": data.Data,
	}
}

func CaptchaStringGet(width, height int) map[string]string {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}

	data, err := defaultCaptcha.StringCaptcha(width, height, 4)
	if err != nil {
		return map[string]string{
			"id":          "",
			"base64image": "",
		}
	}

	return map[string]string{
		"id":          data.CaptchaId,
		"base64image": data.Data,
	}
}

func CaptchaAudioGet(audioStr string) map[string]string {
	data, err := defaultCaptcha.AudioCaptcha(audioStr, 4)
	if err != nil {
		return map[string]string{
			"id":          "",
			"base64image": "",
		}
	}

	return map[string]string{
		"id":          data.CaptchaId,
		"base64image": data.Data,
	}
}

// 验证图片验证码
func CaptchaVerify(id, code string) bool {
	return defaultCaptcha.CaptchaVerify(id, code)
}
