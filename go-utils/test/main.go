package main

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gif-gif/go.io/goio"
	"github.com/gogf/gf/util/gconv"
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
	time.Sleep(30 * time.Second)
}

func testSign() {
	ts := time.Now().Unix()
	sign := goutils.Md5([]byte(gconv.String(ts) + "123456"))
	golog.WithTag("sign").Info(ts, sign)
}
