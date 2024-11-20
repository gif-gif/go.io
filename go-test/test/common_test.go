package main

import (
	"fmt"
	goutils "github.com/gif-gif/go.io/go-utils"
	goxlsx "github.com/gif-gif/go.io/go-xlsx"
	"testing"
	"time"
)

// readCSVAndConvertToJSON 读取 CSV 文件并转换为 JSON 数据

func main() {
	filePath := "/Users/Jerry/Documents/my/test/data/detail.csv" // 替换为你的 CSV 文件路径
	csv, err := goxlsx.NewCsvReader(filePath, '\t')
	if err != nil {
		return
	}

	csv.ReadLineJson(goxlsx.UTF16, func(record map[string]string) error {
		fmt.Println(record)
		return nil
	})

	//jsonData, err := csv.ReadAllJson(goxlsx.UTF16)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}

	// 打印 JSON 数据
	//fmt.Println(jsonData)
}

func TestCsvRead1(t *testing.T) {
	tt, err := goutils.ConvertToGMTTime("2021-08-06T07:00:00+0000")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 将时间对象转换为 UTC 时区
	utcTime := tt.UTC()

	// 格式化为 GMT 时间字符串
	gmtTimeStr := utcTime.Format(time.DateOnly)

	fmt.Println("当前时间的 GMT 时间:", gmtTimeStr, utcTime.Hour())
}
