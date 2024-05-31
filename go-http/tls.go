package gohttp

import (
	"crypto/tls"
	golog "github.com/gif-gif/go.io/go-log"
	"io/ioutil"
)

type Tls struct {
	CaCrtFile     string
	ClientCrtFile string
	ClientKeyFile string
}

func (s *Tls) CaCrt() []byte {
	if s.CaCrtFile == "" {
		return caCert
	}
	bts, err := ioutil.ReadFile(s.CaCrtFile)
	if err != nil {
		golog.Error(err.Error())
	}
	return bts
}

func (s *Tls) ClientCrt() tls.Certificate {
	crt, err := tls.LoadX509KeyPair(s.ClientCrtFile, s.ClientKeyFile)
	if err != nil {
		golog.Error(err.Error())
	}
	return crt
}
