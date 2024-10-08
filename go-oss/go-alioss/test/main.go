package main

import (
	"flag"
	"fmt"
	goalioss "github.com/gif-gif/go.io/go-oss/go-alioss"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio"
	"io/ioutil"
	"os"
	"strings"
)

var (
	AccessKeyId     = flag.String("access_key_id", "", "")
	AccessKeySecret = flag.String("access_key_secret", "", "")
	Endpoint        = flag.String("endpoint", "", "")
	Bucket          = flag.String("bucket", "", "")
	Domain          = flag.String("domain", "", "")
)

func main() {
	goio.FlagInit()

	args := os.Args
	if l := len(args); l < 2 {
		fmt.Println("请选择上传文件!")
		return
	}

	conf := goalioss.Config{
		AccessKeyId:     *AccessKeyId,
		AccessKeySecret: *AccessKeySecret,
		Endpoint:        *Endpoint,
		Bucket:          *Bucket,
		Domain:          *Domain,
	}

	if conf.AccessKeyId == "" {
		conf.AccessKeyId = ""
	}
	if conf.AccessKeySecret == "" {
		conf.AccessKeySecret = ""
	}
	if conf.Endpoint == "" {
		conf.Endpoint = "oss-cn-beijing.aliyuncs.com"
	}
	if conf.Bucket == "" {
		conf.Bucket = ""
	}

	err := goalioss.Init(conf)
	if err != nil {
		return
	}

	for n, i := range args {
		if n == 0 {
			continue
		}

		b, err := ioutil.ReadFile(i)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		index := strings.LastIndex(i, "/")

		var filename string
		if index == -1 {
			filename = i
		} else {
			filename = i[index+1:]
		}

		md5 := goutils.MD5(b)
		filename = fmt.Sprintf("%s/%s/%s", md5[0:2], md5[2:4], filename)

		url, err := goalioss.Default().Upload(strings.ToLower(filename), b)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(url)
	}
}
