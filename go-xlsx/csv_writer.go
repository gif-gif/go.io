package goxlsx

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
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
	writer *csv.Writer //for local file writer

	titles []string   //for http writer
	rows   [][]string //for http writer
	encoding.Encoding
}

func NewCsvWriter(csvFile string, comma rune, encoding ...encoding.Encoding) (*CsvWriter, error) {
	encode := UTF8
	if len(encoding) > 0 {
		encode = encoding[0]
	}
	return &CsvWriter{
		FilePath: csvFile,
		Comma:    comma,
		Encoding: encode,
	}, nil
}

// 不需要创建本地文件， 只写入数据。如：网络IO
func NewCsvWriterNoneFile(comma rune) (*CsvWriter, error) {
	return NewCsvWriter("", comma)
}

func (c *CsvWriter) getWriter() (*csv.Writer, error) {
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

	err = c.WriteBom(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	c.file = file
	// 创建一个 CSV 写入器
	writer := csv.NewWriter(transform.NewWriter(file, c.Encoding.NewEncoder()))
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

	defer file.Close()
	err = c.WriteBom(file)
	if err != nil {
		return err
	}
	// 创建一个 CSV 写入器
	writer := csv.NewWriter(transform.NewWriter(file, c.Encoding.NewEncoder()))
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
// 适合一次性追加多条数据，每次会打开文件，写入，关闭文件，效率低
func (c *CsvWriter) AppendToCSV(data [][]string) error {
	// Create CSV writer
	// 打开 CSV 文件
	// Open file in append mode
	file, err := os.OpenFile(c.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	// 创建一个 CSV 写入器
	writer := csv.NewWriter(transform.NewWriter(file, c.Encoding.NewEncoder()))
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

// 逐行写入, 需要手动调用 c.Close() 或者 end 为true 时才会关闭文件
func (c *CsvWriter) WriteLine(record []string, end ...bool) error {
	if c.writer == nil { //不存在 则创建
		w, err := c.getWriter()
		if err != nil {
			return err
		}
		c.writer = w
		// 写入记录到 CSV 文件
		if err := c.writer.Write(c.titles); err != nil { //如果有titles 则自动写入
			return err
		}
	}

	if end != nil && len(end) > 0 && end[0] {
		defer c.Close()
	}

	// 写入记录到 CSV 文件
	defer c.writer.Flush() // 确保在函数结束时刷新写入器
	// 写入记录到 CSV 文件
	if err := c.writer.Write(record); err != nil {
		return err
	}
	return nil
}

// 关闭文件	c.file.Close()
func (c *CsvWriter) Close() {
	if c.file != nil {
		c.file.Close()
	}
	c.writer = nil
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

// 一次性把titles、rows输出到网络
func (x *CsvWriter) OutputForGin(ctx *gin.Context, filename string) (err error) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	return x.Output(ctx.Writer)
}

// 一次性把titles、rows输出到网络
func (x *CsvWriter) OutputResponseWriter(w http.ResponseWriter, filename string) (err error) {
	header := w.Header()
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Disposition", "attachment; filename="+filename)
	header.Set("Content-Transfer-Encoding", "binary")
	header.Set("Expires", "0")
	return x.Output(w)
}

func (x *CsvWriter) Output(w io.Writer) (err error) {
	err = x.WriteBom(w)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

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

// 写入BOM
func (x *CsvWriter) WriteBom(w io.Writer) error {
	if x.Encoding == UTF8 {
		_, err := w.Write(BOM_UTF8)
		if err != nil {
			return err
		}
	} else if x.Encoding == UTF16 {
		_, err := w.Write(BOM_UTF16)
		if err != nil {
			return err
		}
	}
	return nil
}
