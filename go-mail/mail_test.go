package gomail

import "testing"

func TestMail_Send_126(t *testing.T) {
	conf := Config{
		Username: "goio@126.com",
		Password: "XYQNHPYFSLHYVRGCA",
		Host:     "smtp.126.com",
		Port:     465,
		TLS:      true,
	}

	msg := Message{
		Sender:     "goio@126.com",
		SenderName: "测试者",
		Receivers:  []string{"goio@126.com"},
		Subject:    "hi",
		Body:       "hi goio",
	}

	New(conf).Send(msg)
}

func TestMail_Send_qq(t *testing.T) {
	conf := Config{
		Username: "service@shuzhuo.cn",
		Password: "k27Cicaftj9Nqp5da",
		Host:     "smtp.exmail.qq.com",
		Port:     465,
		TLS:      true,
	}

	msg := Message{
		Sender:     "service@shuzhuo.cn",
		SenderName: "测试者",
		Receivers:  []string{"goio@126.com"},
		Subject:    "hi",
		Body:       `<h3>hi goio</h3><p>请点击下面的链接进行修改密码：http://www.baidu.com</p>`,
	}

	New(conf).Send(msg)
}

func TestMail_Send_gmail(t *testing.T) {
	conf := Config{
		Username: "goio@gmail.com",
		Password: "ljtgepqpmraiixeea",
		Host:     "smtp.gmail.com",
		Port:     465,
		TLS:      true,
	}

	msg := Message{
		Sender:     "goio@gmail.com",
		SenderName: "测试者",
		Receivers:  []string{"goio@126.com"},
		Subject:    "hi",
		Body:       `<h3>hi goio</h3><p>请点击下面的链接进行修改密码：http://www.baidu.com</p>`,
	}

	New(conf).Send(msg)
}
