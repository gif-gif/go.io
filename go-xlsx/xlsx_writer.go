package goxlsx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding"
	"io"
	"net/http"
)

// 扩展需要append to file 功能
type XlsxWrite struct {
	fh        *excelize.File
	sheetName string
	titles    *[]string
	rows      []*[]interface{}
	encoding.Encoding
}

func NewWriter(encoding ...encoding.Encoding) *XlsxWrite {
	encode := UTF8
	if len(encoding) > 0 {
		encode = encoding[0]
	}
	return &XlsxWrite{
		fh:        excelize.NewFile(),
		sheetName: "Sheet1",
		Encoding:  encode,
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
	//if x.Encoding == UTF8 {
	//	_, err = w.Write(BOM_UTF8)
	//	if err != nil {
	//		return err
	//	}
	//} else if x.Encoding == UTF16 {
	//	_, err = w.Write(BOM_UTF16)
	//	if err != nil {
	//		return err
	//	}
	//}

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
