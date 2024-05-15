package gominio

import (
	"context"
	"github.com/minio/minio-go/v7"
)

type (
	Oss interface {
		Init(conf Config) *uploader
		Client() *minio.Client
		ContentType(value string) *uploader
		Options(opts minio.PutObjectOptions) *uploader
		Upload(ctx context.Context, objectName string, filePath string) (*minio.UploadInfo, error)
		CreateBucket(ctx context.Context, bucketName string, location string) error
	}

	customOssModel struct {
		oss *uploader
	}
)

func (c *customOssModel) Client() *minio.Client {
	return c.oss.client
}

func NewOssModel(conf Config) Oss {
	o := &customOssModel{}
	o.Init(conf)
	return o
}

func (c *customOssModel) Init(conf Config) *uploader {
	c.oss, _ = newUploader(conf)
	return c.oss
}

func (c *customOssModel) ContentType(value string) *uploader {
	return c.oss.ContentType(value)
}

func (c *customOssModel) Options(opts minio.PutObjectOptions) *uploader {
	return c.oss.Options(opts)
}

func (c *customOssModel) Upload(ctx context.Context, objectName string, filePath string) (*minio.UploadInfo, error) {
	return c.oss.Upload(ctx, objectName, filePath)
}

func (c *customOssModel) CreateBucket(ctx context.Context, bucketName string, location string) error {
	return c.oss.CreateBucket(ctx, bucketName, location)
}
