package goxlsx

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path"
)

type CsvWriter struct {
	FilePath string
	Comma    rune //csv 列分割符

	//只有逐行写入时才保留这个状态
	file   *os.File
	writer *csv.Writer

	titles []string
	rows   [][]string
}

func NewCsvWriter(csvFile string, comma rune) (*CsvWriter, error) {
	return &CsvWriter{
		FilePath: csvFile,
		Comma:    comma,
	}, nil
}

// 不需要创建本地文件， 只写入数据。如：网络IO
func NewCsvWriterNoneFile(comma rune) (*CsvWriter, error) {
	return NewCsvWriter("", comma)
}

func (c *CsvWriter) SetTitles(titles []string) {
	c.titles = titles
}

func (c *CsvWriter) getTitles() []string {
	return c.titles
}

func (x *CsvWriter) AppendData(data []string) *CsvWriter {
	x.rows = append(x.rows, data)
	return x
}

func (x *CsvWriter) AppendRows(data [][]string) *CsvWriter {
	for _, i := range data {
		x.rows = append(x.rows, i)
	}
	return x
}

func (c *CsvWriter) GetWriter() (*csv.Writer, error) {
	// 打开 CSV 文件
	// Open file in append mode
	dirname := path.Dir(c.FilePath)
	if _, err := os.Stat(dirname); err != nil {
		os.MkdirAll(dirname, 0755)
	}

	file, err := os.OpenFile(c.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	c.file = file
	// 创建一个 CSV 写入器
	writer := csv.NewWriter(file)
	writer.Comma = c.Comma // 使用分号作为分隔符
	return writer, nil
}

// 一次性写入
func (c *CsvWriter) WriteData(records [][]string) error {
	// 写入记录到 CSV 文件
	dirname := path.Dir(c.FilePath)
	if _, err := os.Stat(dirname); err != nil {
		os.MkdirAll(dirname, 0755)
	}
	file, err := os.OpenFile(c.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	// 创建一个 CSV 写入器
	writer := csv.NewWriter(file)
	writer.Comma = c.Comma // 使用分号作为分隔符
	defer writer.Flush()   // 确保在函数结束时刷新写入器
	// 写入记录到 CSV 文件
	if err := writer.Write(c.titles); err != nil { //如果有titles 则自动写入
		return err
	}

	err = writer.WriteAll(records)
	if err != nil {
		return err
	}
	return nil
}

// AppendToCSV appends data to an existing CSV file
func (c *CsvWriter) AppendToCSV(data [][]string) error {
	// Create CSV writer
	// 打开 CSV 文件
	// Open file in append mode
	file, err := os.OpenFile(c.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	// 创建一个 CSV 写入器
	writer := csv.NewWriter(file)
	writer.Comma = c.Comma // 使用分号作为分隔符

	defer writer.Flush()

	// Write all data rows
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

//	func TestCsvWriter(t *testing.T) {
//		csvWriter, err := goxlsx.NewCsvWriter("test.csv", ',')
//		if err != nil {
//			golog.Error(err)
//			return
//		}
//		// 准备要写入的记录
//		records := [][]string{
//			{"Alice", "30", "New York"},
//			{"Bob", "25", "Los Angeles"},
//			{"Charlie", "35", "Chicago"},
//		}
//		titles := []string{"Name", "Age", "City"}
//		csvWriter.SetTitles(titles)
//		defer csvWriter.Close()
//		for _, record := range records {
//			err = csvWriter.WriteLine(record)
//			if err != nil {
//				golog.Error(err)
//			}
//		}
//	}
//
// 逐行写入, 需要手动调用 c.Close()
func (c *CsvWriter) WriteLine(record []string) error {
	if c.writer == nil { //不存在 则创建
		w, err := c.GetWriter()
		if err != nil {
			return err
		}
		c.writer = w

		// 写入记录到 CSV 文件
		if err := c.writer.Write(c.titles); err != nil { //如果有titles 则自动写入
			return err
		}
	}
	// 写入记录到 CSV 文件
	defer c.writer.Flush() // 确保在函数结束时刷新写入器
	// 写入记录到 CSV 文件
	if err := c.writer.Write(record); err != nil {
		return err
	}
	return nil
}

func (c *CsvWriter) Close() {
	if c.file != nil {
		c.file.Close()
	}
	c.writer = nil
}

func (x *CsvWriter) OutputForGin(ctx *gin.Context, filename string) (err error) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	// 写入表头
	return x.Output(ctx.Writer)
}

func (x *CsvWriter) OutputResponseWriter(w http.ResponseWriter, filename string) (err error) {
	header := w.Header()
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Disposition", "attachment; filename="+filename)
	header.Set("Content-Transfer-Encoding", "binary")
	header.Set("Expires", "0")
	return x.Output(w)
}

func (x *CsvWriter) Output(w io.Writer) (err error) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// 写入表头
	if len(x.titles) > 0 {
		if err := writer.Write(x.titles); err != nil {
			return err
		}
	}
	// 写入数据行
	for _, row := range x.rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
