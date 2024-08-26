package goalioss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	golog "github.com/gif-gif/go.io/go-log"
)

var __clients = map[string]*Uploader{}

func Init(configs ...Config) (err error) {
	for _, conf := range configs {
		name := conf.Name
		if name == "" {
			name = "default"
		}

		__clients[name], err = create(conf)
		if err != nil {
			return
		}
	}

	return
}

func create(conf Config) (*Uploader, error) {
	o := &Uploader{
		conf:    conf,
		options: []oss.Option{},
	}

	client, err := o.getClient()
	if err != nil {
		golog.Error(err.Error())
		return nil, err
	}

	o.client = client

	bucket, err := o.getBucket()
	if err != nil {
		golog.Error(err.Error())
		return nil, err
	}

	o.bucket = bucket
	return o, nil
}

func New(conf Config) (*Uploader, error) {
	err := Init(conf)
	if err != nil {
		return nil, err
	}
	return GetClient(conf.Name), nil
}

func GetClient(names ...string) *Uploader {
	name := "default"
	if l := len(names); l > 0 {
		name = names[0]
		if name == "" {
			name = "default"
		}
	}
	if cli, ok := __clients[name]; ok {
		return cli
	}
	return nil
}

func Default() *Uploader {
	if cli, ok := __clients["default"]; ok {
		return cli
	}

	if l := len(__clients); l == 1 {
		for _, cli := range __clients {
			return cli
		}
	}

	golog.WithTag("gominio").Error("no default minio client")

	return nil
}

func (g *Uploader) AliClient() *oss.Client {
	return g.client
}
