package goxlsx

import (
	golog "github.com/gif-gif/go.io/go-log"
	"testing"
)

func TestXlsxWrite(t *testing.T) {
	w := NewWriter()
	w.SetSheetName("test")
	w.titles = &[]string{"title1", "title2", "title3"}
	var data [][]interface{}
	for i := 0; i < 10; i++ {
		data = append(data, []interface{}{"1", "2", "3"})
	}
	w.AppendRows(data)
	err := w.Save2File("test.xlsx")
	if err != nil {
		golog.Error(err)
	}
}

func TestCsvRead(t *testing.T) {
	w, err := NewCsvReader("/Users/Jerry/Documents/my/test/data/detail.csv", ',')
	if err != nil {
		golog.Error(err)
		return
	}

	err = w.ReadUTF16Line(func(record []string) {
		golog.WithTag("record").Info(record)
	})
	if err != nil {
		golog.Error(err)
		return
	}
}

func TestCsvRead1(t *testing.T) {
	w, err := NewCsvReader("/Users/Jerry/Documents/my/test/data/all.csv", '\t')
	if err != nil {
		golog.Error(err)
		return
	}

	data, err := w.ReadUTF16All()
	if err != nil {
		golog.Error(err)
		return
	}

	golog.WithTag("data").Info(data)
}
