package goutils

import (
	"net/mail"
	"regexp"
)

func HideEmail(email string) string {
	re := regexp.MustCompile("(?P<name>[^@]+)@(?P<domain>[^@]+\\.[^@]+)")
	matches := re.FindStringSubmatch(email)

	if len(matches) < 3 {
		return email
	}

	// 隐藏用户名的一部分
	name := matches[1]
	hiddenName := ""
	if len(name) > 3 {
		// 前三个字符保持不变，后面的字符替换为星号
		hiddenName = name[:3] + string(make([]rune, len(name)-3, len(name)-3))
	} else {
		hiddenName = string(make([]rune, len(name), len(name))) // 全部替换为星号
	}

	return hiddenName + "@" + matches[2]
}

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
