package goxlsx

import (
	"encoding/csv"
	"fmt"
	gofile "github.com/gif-gif/go.io/go-file"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"os"
)

type CsvRead struct {
	FilePath string
}

func NewCsvReader(csvFile string) (*CsvRead, error) {
	e, err := gofile.Exist(csvFile)
	if err != nil {
		return nil, err
	}
	if !e {
		return nil, fmt.Errorf("file not exist")
	}
	return &CsvRead{
		FilePath: csvFile,
	}, nil
}

func (c *CsvRead) ReadGBKAll() ([][]string, error) {
	// 打开 CSV 文件
	file, err := os.Open(c.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(transform.NewReader(file, simplifiedchinese.GBK.NewDecoder()))
	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return records, nil
}

func (c *CsvRead) ReadUTF8All() ([][]string, error) {
	// 打开 CSV 文件
	file, err := os.Open(c.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return records, nil
}

func (c *CsvRead) ReadUTF16All() ([][]string, error) {
	// 打开 CSV 文件
	file, err := os.Open(c.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 创建 UTF-16 解码器
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewDecoder()
	//utf16Reader := transform.NewReader(file, decoder)

	// 读取所有内容到内存中
	//utf8Data, err := ioutil.ReadAll(utf16Reader)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read file: %v", err)
	//}

	// 创建 CSV 阅读器
	reader := csv.NewReader(transform.NewReader(file, decoder))
	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return records, nil
}

func (c *CsvRead) ReadUTF8Line(lineDataFunc func(record []string)) error {
	// 打开 CSV 文件
	file, err := os.Open(c.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 创建 CSV 阅读器
	reader := csv.NewReader(file)

	// 按行读取文件
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // 文件读取完毕
			}
			return fmt.Errorf("failed to read line: %v", err)
		}
		// 处理每一行数据
		lineDataFunc(record)
	}

	return nil
}

func (c *CsvRead) ReadUTF16Line(lineDataFunc func(record []string)) error {
	// 打开 CSV 文件
	file, err := os.Open(c.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 创建 UTF-16 解码器
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewDecoder()
	//utf16Reader := transform.NewReader(file, decoder)

	// 读取所有内容到内存中
	//utf8Data, err := ioutil.ReadAll(utf16Reader)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read file: %v", err)
	//}

	// 创建 CSV 阅读器
	reader := csv.NewReader(transform.NewReader(file, decoder))

	// 按行读取文件
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // 文件读取完毕
			}
			return fmt.Errorf("failed to read line: %v", err)
		}
		// 处理每一行数据
		lineDataFunc(record)
	}

	return nil
}

func (c *CsvRead) ReadGBKLine(lineDataFunc func(record []string)) error {
	// 打开 CSV 文件
	file, err := os.Open(c.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(transform.NewReader(file, simplifiedchinese.GBK.NewDecoder()))
	// 按行读取文件
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // 文件读取完毕
			}
			return fmt.Errorf("failed to read line: %v", err)
		}
		// 处理每一行数据
		lineDataFunc(record)
	}

	return nil
}
