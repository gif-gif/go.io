package gominio

import (
	"context"
	"github.com/minio/minio-go/v7"
)

type (
	Oss interface {
		Init(conf Config) *GoMinio
		Client() *minio.Client
		FPutObject(ctx context.Context, objectName string, sourceFilePath string, options *minio.PutObjectOptions) (*minio.UploadInfo, error)
		FGetObject(ctx context.Context, objectName, saveFilePath string, options *minio.GetObjectOptions) error
		CreateBucket(ctx context.Context, bucketName string, location string) error
	}

	customOssModel struct {
		oss *GoMinio
	}
)

func NewOssModel(conf Config) Oss {
	o := &customOssModel{}
	o.Init(conf)
	return o
}

func (c *customOssModel) Client() *minio.Client {
	return c.oss.client
}

func (c *customOssModel) Init(conf Config) *GoMinio {
	c.oss, _ = newGoMinio(conf)
	return c.oss
}

func (c *customOssModel) FPutObject(ctx context.Context, objectName string, sourceFilePath string, options *minio.PutObjectOptions) (*minio.UploadInfo, error) {
	return c.oss.FPutObject(ctx, objectName, sourceFilePath, options)
}

func (c *customOssModel) FGetObject(ctx context.Context, objectName, saveFilePath string, options *minio.GetObjectOptions) error {
	return c.oss.FGetObject(ctx, objectName, saveFilePath, options)
}

func (c *customOssModel) CreateBucket(ctx context.Context, bucketName string, location string) error {
	return c.oss.CreateBucket(ctx, bucketName, location)
}
