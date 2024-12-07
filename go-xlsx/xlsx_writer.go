package goxlsx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
)

// 扩展需要append to file 功能
type XlsxWrite struct {
	fh        *excelize.File
	sheetName string
	titles    *[]string
	rows      []*[]interface{}
}

func NewWriter() *XlsxWrite {
	return &XlsxWrite{
		fh:        excelize.NewFile(),
		sheetName: "Sheet1",
	}
}

func (x *XlsxWrite) SetTitles(titles []string) *XlsxWrite {
	x.titles = &titles
	return x
}

func (x *XlsxWrite) AppendData(data []interface{}) *XlsxWrite {
	x.rows = append(x.rows, &data)
	return x
}

func (x *XlsxWrite) AppendRows(data [][]interface{}) *XlsxWrite {
	for _, i := range data {
		x.rows = append(x.rows, &i)
	}
	return x
}

func (x *XlsxWrite) SetSheetName(sheetName string) *XlsxWrite {
	x.sheetName = sheetName
	index, err := x.fh.NewSheet(x.sheetName)
	if err != nil {
		return nil
	}
	x.fh.SetActiveSheet(index)
	return x
}

func (x *XlsxWrite) Save2File(filename string) (err error) {
	err = x.fh.SetSheetRow(x.sheetName, "A1", x.titles)
	if err != nil {
		return err
	}

	for i := 0; i < len(x.rows); i++ {
		err = x.fh.SetSheetRow(x.sheetName, fmt.Sprintf("A%d", i+2), x.rows[i])
		if err != nil {
			return err
		}
	}

	if err = x.fh.SaveAs(filename); err != nil {
		return err
	}

	return nil
}

func (x *XlsxWrite) OutputForGin(ctx *gin.Context, filename string) (err error) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Expires", "0")
	return x.Output(ctx.Writer)
}

func (x *XlsxWrite) OutputResponseWriter(w http.ResponseWriter, filename string) (err error) {
	header := w.Header()
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Disposition", "attachment; filename="+filename)
	header.Set("Content-Transfer-Encoding", "binary")
	header.Set("Expires", "0")
	return x.Output(w)
}

func (x *XlsxWrite) Output(w io.Writer) (err error) {
	err = x.fh.SetSheetRow(x.sheetName, "A1", x.titles)
	if err != nil {
		return err
	}

	for i := 0; i < len(x.rows); i++ {
		err = x.fh.SetSheetRow(x.sheetName, fmt.Sprintf("A%d", i+2), x.rows[i])
		if err != nil {
			return err
		}
	}

	if err = x.fh.Write(w); err != nil {
		return err
	}

	return nil
}
