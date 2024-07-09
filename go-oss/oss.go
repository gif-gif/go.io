package gooss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type GoOss struct {
	oss *Uploader
}

func Init(conf Config) *GoOss {
	__oss, _ := New(conf)
	return &GoOss{oss: __oss}
}

func (g *GoOss) Client() *oss.Client {
	return g.oss.client
}

func (g *GoOss) ContentType(value string) *Uploader {
	return g.oss.ContentType(value)
}

func (g *GoOss) Options(opts ...oss.Option) *Uploader {
	return g.oss.Options(opts...)
}

func (g *GoOss) Upload(filename string, body []byte) (string, error) {
	return g.oss.Upload(filename, body)
}
