package main

import (
	"fmt"
	goxlsx "github.com/gif-gif/go.io/go-xlsx"
)

// readCSVAndConvertToJSON 读取 CSV 文件并转换为 JSON 数据

func main() {
	filePath := "/Users/Jerry/Documents/my/test/data/all.csv" // 替换为你的 CSV 文件路径
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
