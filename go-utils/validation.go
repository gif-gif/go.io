package goutils

import "regexp"

func ValidPassword(str string) (msg string, matched bool) {
	msg = "至少一位数字、大小字母,且长度6-20位"

	matched, _ = regexp.MatchString("[0-9]+", str)
	if !matched {
		return
	}

	matched, _ = regexp.MatchString("[a-z]+", str)
	if !matched {
		return
	}

	matched, _ = regexp.MatchString("[A-Z]+", str)
	if !matched {
		return
	}

	if l := len(str); l < 6 || l > 20 {
		matched = false
		return
	}

	msg = ""
	matched = true

	return
}

func ValidPasswordV2(str string) (msg string, matched bool) {
	msg = "至少一位数字、大小字母和特殊字符,且长度6-20位"

	matched, _ = regexp.MatchString("[0-9]+", str)
	if !matched {
		return
	}

	matched, _ = regexp.MatchString("[a-z]+", str)
	if !matched {
		return
	}

	matched, _ = regexp.MatchString("[A-Z]+", str)
	if !matched {
		return
	}

	matched, _ = regexp.MatchString("[^0-9a-zA-Z]+", str)
	if !matched {
		return
	}

	if l := len(str); l < 6 || l > 20 {
		matched = false
		return
	}

	msg = ""
	matched = true

	return
}
