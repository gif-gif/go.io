package gominio

import (
	"context"
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/url"
	"time"
)

type Uploader struct {
	conf   Config
	client *minio.Client
}

func Create(conf Config) (*Uploader, error) {
	o := &Uploader{
		conf: conf,
	}

	client, err := o.getClient()
	if err != nil {
		golog.Error(err.Error())
		return nil, err
	}

	o.client = client

	return o, nil
}

func (o *Uploader) FPutObject(ctx context.Context, objectName, filePath string, options *minio.PutObjectOptions) (*minio.UploadInfo, error) {
	if objectName == "" {
		return nil, errors.New("文件名为空")
	}

	if options == nil {
		options = &minio.PutObjectOptions{ContentType: "application/octet-stream"}
	}

	info, err := o.client.FPutObject(ctx, o.conf.Bucket, objectName, filePath, *options)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (o *Uploader) FGetObject(ctx context.Context, objectName, saveFilePath string, options *minio.GetObjectOptions) error {
	if objectName == "" {
		return errors.New("文件名为空")
	}

	if options == nil {
		options = &minio.GetObjectOptions{}
	}

	err := o.client.FGetObject(ctx, o.conf.Bucket, objectName, saveFilePath, *options)
	if err != nil {
		return err
	}

	return nil
}

// UploadExpiredObject
// UploadStreamObject

func (o *Uploader) CreateBucket(ctx context.Context, bucketName string, location string) error {
	err := o.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := o.client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (o *Uploader) ListBuckets() ([]minio.BucketInfo, error) {
	return o.client.ListBuckets(context.Background())
}

func (o *Uploader) BucketExists(bucketName string) (bool, error) {
	return o.client.BucketExists(context.Background(), bucketName)
}

func (o *Uploader) ListObjects(bucketName string, opts *minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	if opts == nil {
		opts = &minio.ListObjectsOptions{}
	}
	return o.client.ListObjects(context.Background(), bucketName, *opts)
}

func (o *Uploader) PresignedGetObject(bucketName, objectName string, expiry time.Duration, reqParams url.Values) (u *url.URL, err error) {
	return o.client.PresignedGetObject(context.Background(), bucketName, objectName, expiry, reqParams)
}

func (o *Uploader) getClient() (*minio.Client, error) {
	minioClient, err := minio.New(o.conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(o.conf.AccessKeyId, o.conf.AccessKeySecret, ""),
		Secure: o.conf.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient, nil
}
