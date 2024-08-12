package test

import (
	"context"
	"flag"
	"fmt"
	gominio "github.com/gif-gif/go.io/go-oss/go-minio"
	"github.com/minio/minio-go/v7"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	AccessKeyId     = flag.String("access_key_id", "fPztco0yxWC1qrxz6iDN", "")
	AccessKeySecret = flag.String("access_key_secret", "kWS6q1PIxjcEnLSkoBWrWqPnsHFHm1Q8cvG1CUPm", "")
	Endpoint        = flag.String("endpoint", "minio.gif00.com", "")
	Bucket          = flag.String("bucket", "test", "")
	Domain          = flag.String("domain", "", "")
)

func createBucketTest() {
	conf := gominio.Config{
		AccessKeyId:     *AccessKeyId,
		AccessKeySecret: *AccessKeySecret,
		Endpoint:        *Endpoint,
		Bucket:          *Bucket,
		Domain:          *Domain,
		UseSSL:          false,
	}

	oss, _ := gominio.Create(conf)
	bucketName := "testbucket"
	location := "us-east-1"

	err := oss.CreateBucket(context.Background(), bucketName, location)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Created bucket successfully")
	}

}

func uploadTest() {
	args := os.Args
	if l := len(args); l < 2 {
		fmt.Println("请选择上传文件!")
		return
	}

	conf := gominio.Config{
		AccessKeyId:     *AccessKeyId,
		AccessKeySecret: *AccessKeySecret,
		Endpoint:        *Endpoint,
		Bucket:          *Bucket,
		Domain:          *Domain,
		UseSSL:          false,
	}

	oss, _ := gominio.Create(conf)

	for n, i := range args {
		if n == 0 {
			continue
		}

		var filename string
		index := strings.LastIndex(i, "/")
		if index == -1 {
			filename = i
		} else {
			filename = i[index+1:]
		}

		info, err := oss.FPutObject(context.Background(), strings.ToLower(filename), i, &minio.PutObjectOptions{
			Expires: time.Now().AddDate(0, 0, 1),
		})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(info)
	}
}

func getTest() {

	conf := gominio.Config{
		AccessKeyId:     *AccessKeyId,
		AccessKeySecret: *AccessKeySecret,
		Endpoint:        *Endpoint,
		Bucket:          *Bucket,
		Domain:          *Domain,
		UseSSL:          false,
	}

	err := gominio.Init(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = gominio.Default().Client().FGetObject(context.Background(), "/test/2024/05/422744271b108960a4818cc91a1822d9.log", "/Users/Jerry/Desktop/bak202405/422744271b108960a4818cc91a1822d9.log", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println("get object succeeded")
	}

	for object := range gominio.Default().Client().ListObjects("test", nil) {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println("found:", object)
	}

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err := gominio.Default().Client().PresignedGetObject("test", "test.apk", time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)
}

func main() {
	//uploadTest()
	//createBucketTest()
	getTest()
}
