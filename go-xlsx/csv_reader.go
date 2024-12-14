package goxlsx

import (
	"encoding/csv"
	"fmt"
	gofile "github.com/gif-gif/go.io/go-file"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"os"
)

type CsvReader struct {
	FilePath string
	Comma    rune //csv 列分割符
	file     *os.File
}

// comma 默认传 ','
func NewCsvReader(csvFile string, comma rune) (*CsvReader, error) {
	e, err := gofile.Exist(csvFile)
	if err != nil {
		return nil, err
	}
	if !e {
		return nil, fmt.Errorf("file not exist")
	}
	return &CsvReader{
		FilePath: csvFile,
		Comma:    comma,
	}, nil
}

func (c *CsvReader) getReader(encoding encoding.Encoding) (*csv.Reader, error) {
	// 打开 CSV 文件
	file, err := os.Open(c.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	c.file = file

	reader := csv.NewReader(transform.NewReader(file, encoding.NewDecoder()))
	reader.Comma = c.Comma

	return reader, nil
}

func (c *CsvReader) ReadGBKAll() ([][]string, error) {
	// 打开 CSV 文件
	return c.ReadAll(GBK)
}

func (c *CsvReader) ReadUTF8All() ([][]string, error) {
	// 打开 CSV 文件
	return c.ReadAll(UTF8)
}

func (c *CsvReader) ReadUTF16All() ([][]string, error) {
	return c.ReadAll(UTF16)
}

func (c *CsvReader) ReadUTF8Line(lineDataFunc func(record []string) error) error {
	return c.ReadLine(UTF8, lineDataFunc)
}

func (c *CsvReader) ReadUTF16Line(lineDataFunc func(record []string) error) error {
	// 创建 UTF-16 解码器
	return c.ReadLine(UTF16, lineDataFunc)
}

func (c *CsvReader) ReadGBKLine(lineDataFunc func(record []string) error) error {
	return c.ReadLine(GBK, lineDataFunc)
}

func (c *CsvReader) ReadAll(encoding encoding.Encoding) ([][]string, error) {
	// 打开 CSV 文件
	// 创建 UTF-16 解码器
	reader, err := c.getReader(encoding)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return records, nil
}

func (c *CsvReader) ReadLine(encoding encoding.Encoding, lineDataFunc func(record []string) error) error {
	// 打开 CSV 文件
	reader, err := c.getReader(encoding)
	if err != nil {
		return err
	}
	defer c.Close()
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
		e := lineDataFunc(record)
		if e != nil {
			return e
		}
	}
	return nil
}

// 第一行作为字段名称，后续行数据转换为json数据,一行回调一个json数据
func (c *CsvReader) ReadLineJson(encoding encoding.Encoding, lineDataFunc func(record map[string]string) error) error {
	// 打开 CSV 文件
	reader, err := c.getReader(encoding)
	if err != nil {
		return err
	}
	defer c.Close()
	// 获取标题行
	// 创建一个切片来保存 JSON 对象
	// 遍历每一行（从第二行开始，因为第一行是标题）
	record, err := reader.Read()
	if err != nil {
		return err
	}
	headers := record
	// 遍历每一行（从第二行开始，因为第一行是标题）
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
		row := make(map[string]string)
		for i, value := range record {
			row[headers[i]] = value
		}

		e := lineDataFunc(row)
		if e != nil {
			return e
		}
	}
	return nil
}

// 第一行作为字段名称，后续行数据转换为json数据一次性返回所有数据
func (c *CsvReader) ReadAllJson(encoding encoding.Encoding) ([]map[string]string, error) {
	// 打开 CSV 文件
	reader, err := c.getReader(encoding)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	// 读取 CSV 文件的所有内容
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %v", err)
	}

	// 检查 CSV 文件是否为空
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file does not contain enough data")
	}

	// 获取标题行
	headers := records[0]

	// 创建一个切片来保存 JSON 对象
	var jsonData []map[string]string

	// 遍历每一行（从第二行开始，因为第一行是标题）
	for _, record := range records[1:] {
		// 创建一个 map 来保存每一行的数据
		row := make(map[string]string)
		for i, value := range record {
			row[headers[i]] = value
		}
		// 将 map 添加到 jsonData 切片中
		jsonData = append(jsonData, row)
	}

	return jsonData, nil
}

func (c *CsvReader) ReadTitles(encoding encoding.Encoding) ([]string, error) {
	// 打开 CSV 文件
	reader, err := c.getReader(encoding)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	// 获取标题行
	// 第一行是标题视为 标题
	record, err := reader.Read()
	if err != nil {
		return nil, err
	}
	headers := record
	return headers, nil
}

func (c *CsvReader) Close() {
	if c.file != nil {
		c.file.Close()
	}
}
