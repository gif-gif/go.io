package goxlsx

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

// 需要扩展按行读取文件
type XlsxRead struct {
	XlsxFile *excelize.File
}

func NewReader(xlsxFile string) (*XlsxRead, error) {
	f, err := excelize.OpenFile(xlsxFile)
	if err != nil {
		return nil, err
	}

	return &XlsxRead{
		XlsxFile: f,
	}, nil
}

func (r *XlsxRead) GetRows(sheetName string) error {
	rows, err := r.XlsxFile.GetRows(sheetName)
	if err != nil {
		return err
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
	}

	return nil
}

func (r *XlsxRead) Close() error {
	return r.XlsxFile.Close()
}
