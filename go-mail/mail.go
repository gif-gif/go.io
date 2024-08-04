package gomail

import (
	"crypto/tls"
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"io"
	"net"
	"net/smtp"
)

type iMail interface {
	Send(msg Message) error
}

type mail struct {
	conf Config
}

func New(conf Config) iMail {
	return &mail{
		conf: conf,
	}
}

func (m *mail) Send(msg Message) (err error) {
	if msg.Sender == "" {
		msg.Sender = m.conf.Username
	}

	var (
		conn net.Conn
		cli  *smtp.Client
	)

	conn, cli, err = m.client()
	if err != nil {
		return err
	}
	defer conn.Close()
	defer cli.Close()

	if err = cli.Auth(m.auth()); err != nil {
		golog.Error(err.Error())
		return
	}

	if err = cli.Mail(msg.Sender); err != nil {
		golog.Error(err.Error())
		return
	}

	for _, receiver := range msg.Receivers {
		if err = cli.Rcpt(receiver); err != nil {
			golog.Error(err.Error())
			return
		}
	}

	var (
		w io.WriteCloser
	)

	if w, err = cli.Data(); err != nil {
		golog.Error(err.Error())
		return
	}
	defer w.Close()

	if _, err = w.Write(msg.Html()); err != nil {
		golog.Error(err.Error())
		return
	}

	cli.Quit()

	return
}

func (m *mail) client() (conn net.Conn, cli *smtp.Client, err error) {
	addr := fmt.Sprintf("%s:%d", m.conf.Host, m.conf.Port)

	if m.conf.TLS {
		conn, err = tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = net.Dial("tcp", addr)
	}

	if err != nil {
		golog.Error(err.Error())
		return
	}

	cli, err = smtp.NewClient(conn, m.conf.Host)
	if err != nil {
		golog.Error(err.Error())
		return
	}

	return
}

func (m *mail) auth() (auth smtp.Auth) {
	return smtp.PlainAuth("", m.conf.Username, m.conf.Password, m.conf.Host)
}
