package main

import (
	"fmt"
	gocontext "github.com/gif-gif/go.io/go-context"
	gofile "github.com/gif-gif/go.io/go-file"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio"
	"github.com/gogf/gf/util/gconv"
	"sort"
	"strings"
	"time"
)

// PayPlanItemConfig 结构体定义
type PayPlanItemConfig struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Price  float64 `json:"price"`
	IRT    int64   `json:"IRT,optional"`
	INR    int64   `json:"INR,optional"`
	MMK    int64   `json:"mmk,optional"`
	Weight int64   `json:"weight,optional"`
}

func main() {
	goio.Init(goio.DEVELOPMENT)
	//testGenerateAesKeys()
	testSha1Sign()
	<-gocontext.Cancel().Done()
}

func testGetFieldValue() {
	config := PayPlanItemConfig{
		Id:     "123",
		Title:  "Sample Plan",
		Price:  99.99,
		IRT:    1000,
		INR:    2000,
		MMK:    3000,
		Weight: 50,
	}

	// 获取字段值
	fieldName := "Price"
	value, err := goutils.GetFieldValue(&config, fieldName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("The value of field '%s' is: %v\n", fieldName, value)
}

func testSign() {
	ts := time.Now().Unix()
	sign := goutils.Md5([]byte(gconv.String(ts) + "123456"))
	golog.WithTag("sign").Info(ts, sign)
}

func testSha1Sign() {

	///api/mp/v1/verify?signature=950ef40c2493f20071382c888ee0858172a3b08e&echostr=4193083647034290167&timestamp=1725451984&nonce=138798064
	signature := "950ef40c2493f20071382c888ee0858172a3b08e"
	//echoStr := "4193083647034290167"
	nonce := "1387980644"
	var timestamp int64 = 1725451984
	secret := "123qwe"
	args := []string{secret, gconv.String(timestamp), nonce}
	sort.Strings(args)
	golog.WithTag("sign").Info(strings.Join(args, ""))

	sign := goutils.CheckSignSha1(secret, nonce, 20000, timestamp, signature)
	golog.WithTag("sign").Info(sign)
}

func testFileMd5() {
	filePath := "/Users/Jerry/Downloads/chrome/dy16.5.0.apk"
	md5, err := goutils.CalculateFileMD5("/Users/Jerry/Downloads/chrome/dy16.5.0.apk")
	if err != nil {
		golog.WithTag("md5").Error(err)
	} else {
		golog.WithTag("md5").Info(md5)
	}

	body, err := gofile.GetFileContent(filePath)
	if err != nil {
		golog.WithTag("md5").Error(err)
		return
	}

	md5 = goutils.Md5(body)
	golog.WithTag("md5").Info(md5)

	//cea3b6aa0c114de15ba2741e679e91d3
}

func testRaceSpeed() {
	var fns []func()
	fns = append(fns, func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Hello 5")
	})

	fns = append(fns, func() {
		time.Sleep(10 * time.Second)
		fmt.Println("Hello 1")
	})

	fns = append(fns, func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Hello 3")
	})
}

func testGenerateAesKeys() {
	// 加密
	key, iv, err := goutils.GenerateAESKeyAndIV()
	if err != nil {
		golog.WithTag("aes").Error(err)
		return
	}
	golog.WithTag("aesKey").Info(key)
	golog.WithTag("aesIv").Info(iv)
	key, err = goutils.GenerateAESKey()
	if err != nil {
		golog.WithTag("aes").Error(err)
		return
	}
	golog.WithTag("aesKey").Info(key)
}
