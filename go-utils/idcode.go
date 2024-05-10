package goutils

import (
	"bytes"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const key = "6A7CDKV5TBH0ULFSEP82XMW1G9R3YJQNZ"

type idCode struct {
	base int64
	key  string
	l    int
}

func NewIdCode(key string) *idCode {
	return &idCode{
		base: 100000,
		key:  key,
		l:    len(key),
	}
}

/**
 * <随机字符A><加密字符串><验证字符B>
 * 随机字符A = 从key里面随机获取一个字符
 * 加密字符串 = 每位转换成16进制的ID数字+随机数，然后取余key的长度，得到的数字就是该位在key的字符，加入到总加密后的字符串中
 * 验证字符B = 从key里面获取一个字符，字符位置=(密钥长度-随机数+给定ID长度)%密钥长度
 */
func (c *idCode) Encode(id int64) string {
	id += c.base

	n := c.randNum()

	hexArr := []rune(strconv.FormatInt(id, 16))
	keyArr := []rune(c.key)

	var buf bytes.Buffer
	buf.WriteRune(keyArr[n])

	for _, h := range hexArr {
		hi, _ := strconv.ParseInt(string(h), 16, 64)
		offset := int(hi) + n
		buf.WriteRune(keyArr[offset%c.l])
	}

	buf.WriteRune(keyArr[(c.l-n+len(strconv.FormatInt(id, 10)))%c.l])

	return buf.String()
}

/**
 * <随机字符A><加密字符串><验证字符B>
 * 1. 根据随机字符A，获取随机数n
 * 2. 获取中间字符串，计算得到id
 * 3. 验证验证字符B是否正确
 */
func (c *idCode) Decode(str string) (id int64, err error) {
	if str == "" {
		err = errors.New("code为空")
		return
	}

	strArr := []rune(str)
	keyArr := []rune(c.key)

	l := len(strArr)
	n := strings.IndexRune(c.key, strArr[0])

	var buf bytes.Buffer
	for _, s := range strArr[1 : l-1] {
		pos := strings.IndexRune(c.key, s)
		if pos >= n {
			buf.WriteString(strconv.FormatInt(int64(pos-n), 16))
		} else {
			buf.WriteString(strconv.FormatInt(int64(c.l-n+pos), 16))
		}
	}

	id, err = strconv.ParseInt(buf.String(), 16, 64)
	if err != nil {
		return
	}
	if strArr[l-1] != keyArr[(c.l-n+len(strconv.FormatInt(id, 10)))%c.l] {
		err = errors.New("校验错误")
		return
	}

	id -= c.base

	return
}

func (c *idCode) randNum() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(c.l - 1)
}

func Id2Code(id int64) string {
	return NewIdCode(key).Encode(id)
}

func Code2Id(code string) (int64, error) {
	return NewIdCode(key).Decode(code)
}
