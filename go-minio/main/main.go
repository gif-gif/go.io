package main

import (
	"context"
	"flag"
	"fmt"
	gominio "github.com/jiriyao/go.io/go-minio"
	"os"
	"strings"
)

var (
	AccessKeyId     = flag.String("access_key_id", "fPztco0yxWC1qrxz6iDN", "")
	AccessKeySecret = flag.String("access_key_secret", "kWS6q1PIxjcEnLSkoBWrWqPnsHFHm1Q8cvG1CUPm", "")
	Endpoint        = flag.String("endpoint", "minio.gif00.com", "")
	Bucket          = flag.String("bucket", "test", "")
	Domain          = flag.String("domain", "", "")
)

func main() {
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

	oss := gominio.NewOssModel(conf)

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

		info, err := oss.Upload(context.Background(), strings.ToLower(filename), i)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(info)
	}

}
