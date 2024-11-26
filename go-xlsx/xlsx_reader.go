package goxlsx

import (
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/xuri/excelize/v2"
)

// 需要扩展按行读取文件
type XlsxRead struct {
	FilePath string
}

func NewReader(xlsxFile string) (*XlsxRead, error) {
	return &XlsxRead{
		FilePath: xlsxFile,
	}, nil
}

func (r *XlsxRead) ReadBySheet(sheet string, fn func(n int, row []string) error) error {
	xlsx, err := excelize.OpenFile(r.FilePath)
	if err != nil {
		golog.Error(err)
		return err
	}
	defer xlsx.Close()

	if sheet == "" {
		sheet = xlsx.GetSheetName(0)
	}

	rows, err := xlsx.Rows(sheet)
	if err != nil {
		golog.Error(err)
		return err
	}
	defer rows.Close()

	var n int
	for rows.Next() {
		n++
		row, _ := rows.Columns()
		if err = fn(n, row); err != nil {
			return err
		}
	}

	return nil
}

func (r *XlsxRead) Read(fn func(n int, row []string) error) error {
	return r.ReadBySheet("", fn)
}

func (r *XlsxRead) ReadFile(file string, fn func(n int, row []string) error) error {
	return r.Read(fn)
}
