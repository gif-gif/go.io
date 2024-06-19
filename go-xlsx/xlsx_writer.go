package goxlsx

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gin-gonic/gin"
)

func Writer() *xlsxWrite {
	return &xlsxWrite{
		fh:        excelize.NewFile(),
		sheetName: "Sheet1",
	}
}

type xlsxWrite struct {
	fh        *excelize.File
	sheetName string
	titles    *[]string
	rows      []*[]interface{}
}

func (x *xlsxWrite) SetTitles(titles []string) *xlsxWrite {
	x.titles = &titles
	return x
}

func (x *xlsxWrite) SetData(data []interface{}) *xlsxWrite {
	x.rows = append(x.rows, &data)
	return x
}

func (x *xlsxWrite) SetRows(data [][]interface{}) *xlsxWrite {
	for _, i := range data {
		x.rows = append(x.rows, &i)
	}
	return x
}

func (x *xlsxWrite) SetSheetName(sheetName string) *xlsxWrite {
	x.sheetName = sheetName
	return x
}

func (x *xlsxWrite) Save2File(filename string) (err error) {
	x.fh.SetSheetRow(x.sheetName, "A1", x.titles)

	for i := 0; i < len(x.rows); i++ {
		x.fh.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), x.rows[i])
	}

	if err = x.fh.SaveAs(filename); err != nil {
		golog.Error(err)
		return
	}

	return nil
}

func (x *xlsxWrite) Output(ctx *gin.Context, filename string) (err error) {
	x.fh.SetSheetRow(x.sheetName, "A1", x.titles)

	for i := 0; i < len(x.rows); i++ {
		x.fh.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), x.rows[i])
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err = x.fh.Write(ctx.Writer); err != nil {
		golog.Error(err)
		return
	}

	return
}
