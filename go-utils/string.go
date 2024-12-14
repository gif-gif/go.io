package goutils

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
	"unicode"
)

// 多字符切割，默认支持逗号，分号，\n
func Split(s string, rs ...rune) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		for _, rr := range rs {
			if rr == r {
				return true
			}
		}
		return r == ',' || r == '，' || r == ';' || r == '；' || r == '\n'
	})
}

// 驼峰转下划线
func Camel2Case(str string) string {
	var bf bytes.Buffer

	for i, r := range str {
		if !unicode.IsUpper(r) {
			bf.WriteRune(r)
			continue
		}
		if i > 0 {
			bf.WriteString("_")
		}
		bf.WriteRune(unicode.ToLower(r))
	}

	return bf.String()
}

// 下划线转驼峰
func Case2Camel(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}

// 如果只需要转换单个字符串为 GBK
func UTF8ToGBK(text string) (string, error) {
	encoder := simplifiedchinese.GBK.NewEncoder()
	gbkBytes, err := encoder.Bytes([]byte(text))
	if err != nil {
		return "", err
	}
	return string(gbkBytes), nil
}

// GBK 转 UTF8
func GBKToUTF8(text string) (string, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	bytes, err := decoder.Bytes([]byte(text))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
