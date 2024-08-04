package goo_mail

import (
	"bytes"
	"fmt"
	"strings"
)

type Message struct {
	Sender     string   // 发送者
	SenderName string   // 发送者名称
	Receivers  []string // 接受者
	Subject    string   // 主题
	Body       string   // 内容
}

func (m *Message) Html() []byte {
	var bf bytes.Buffer

	bf.WriteString("To: " + strings.Join(m.Receivers, ",\r\n "))
	bf.WriteString("\r\n")
	if m.SenderName != "" {
		bf.WriteString(fmt.Sprintf("From: \"%s\"<%s>", m.SenderName, m.Sender))
	} else {
		bf.WriteString("From: " + m.Sender)
	}
	bf.WriteString("\r\n")
	bf.WriteString("Subject: " + m.Subject)
	bf.WriteString("\r\n")
	bf.WriteString("Content-Type:text/html;charset=utf-8")
	bf.WriteString("\r\n\r\n")
	bf.WriteString(m.Body)

	return bf.Bytes()
}
