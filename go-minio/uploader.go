package gominio

import (
	"context"
	"errors"
	golog "github.com/jiriyao/go.io/go-log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type uploader struct {
	conf    Config
	client  *minio.Client
	options minio.PutObjectOptions
}

func newUploader(conf Config) (*uploader, error) {
	o := &uploader{
		conf:    conf,
		options: minio.PutObjectOptions{},
	}

	client, err := o.getClient()
	if err != nil {
		golog.Error(err.Error())
		return nil, err
	}

	o.client = client

	return o, nil
}

func (o *uploader) ContentType(value string) *uploader {
	o.options.ContentType = value
	return o
}

func (o *uploader) Options(opts minio.PutObjectOptions) *uploader {
	o.options = opts
	return o
}

func (o *uploader) Upload(ctx context.Context, objectName, filePath string) (*minio.UploadInfo, error) {
	if objectName == "" {
		return nil, errors.New("文件名为空")
	}

	info, err := o.client.FPutObject(ctx, o.conf.Bucket, objectName, filePath, o.options)
	if err != nil {
		log.Fatalln(err)
	}

	return &info, nil
}

// UploadExpiredObject
// UploadStreamObject

func (o *uploader) CreateBucket(ctx context.Context, bucketName string, location string) error {
	err := o.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := o.client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	return nil
}

func (o *uploader) getClient() (*minio.Client, error) {
	minioClient, err := minio.New(o.conf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(o.conf.AccessKeyId, o.conf.AccessKeySecret, ""),
		Secure: o.conf.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return minioClient, nil
}
