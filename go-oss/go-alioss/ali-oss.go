package goalioss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type GoAliOss struct {
	oss *Uploader
}

func Init(conf Config) *GoAliOss {
	__oss, _ := New(conf)
	return &GoAliOss{oss: __oss}
}

func (g *GoAliOss) Client() *oss.Client {
	return g.oss.client
}

func (g *GoAliOss) ContentType(value string) *Uploader {
	return g.oss.ContentType(value)
}

func (g *GoAliOss) Options(opts ...oss.Option) *Uploader {
	return g.oss.Options(opts...)
}

func (g *GoAliOss) Upload(filename string, body []byte) (string, error) {
	return g.oss.Upload(filename, body)
}
