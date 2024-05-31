package gooss

import (
	"bytes"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"path"
	"strings"
)

type uploader struct {
	conf    Config
	client  *oss.Client
	bucket  *oss.Bucket
	options []oss.Option
}

func New(conf Config) (*uploader, error) {
	o := &uploader{
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

func (o *uploader) ContentType(value string) *uploader {
	o.options = append(o.options, oss.ContentType(value))
	return o
}

func (o *uploader) Options(opts ...oss.Option) *uploader {
	o.options = append(o.options, opts...)
	return o
}

func (o *uploader) Upload(filename string, body []byte) (string, error) {
	if filename == "" {
		return "", errors.New("文件名为空")
	}

	md5str := goutils.MD5(body)

	ext := path.Ext(filename)
	index := strings.Index(filename, ext)
	filename = filename[:index] + "_" + md5str[8:24] + filename[index:]

	if err := o.bucket.PutObject(filename, bytes.NewReader(body), o.options...); err != nil {
		golog.Error(err.Error())
		return "", err
	}

	if filename[0:1] != "/" {
		filename = "/" + filename
	}

	if o.conf.Domain != "" {
		if idx, l := strings.LastIndex(o.conf.Domain, "/"), len(o.conf.Domain); idx+1 == l {
			o.conf.Domain = o.conf.Domain[:l-1]
		}
		return o.conf.Domain + filename, nil
	}

	url := "https://" + o.conf.Bucket + "." + o.conf.Endpoint + filename
	return url, nil
}

func (o *uploader) getClient() (*oss.Client, error) {
	return oss.New(o.conf.Endpoint, o.conf.AccessKeyId, o.conf.AccessKeySecret)
}

func (o *uploader) getBucket() (*oss.Bucket, error) {
	return o.client.Bucket(o.conf.Bucket)
}
