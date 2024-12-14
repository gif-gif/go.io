package goxlsx

import (
	golog "github.com/gif-gif/go.io/go-log"
	"testing"
)

func TestXlsxWrite(t *testing.T) {
	w := NewWriter()
	w.SetSheetName("test")
	w.titles = &[]string{"测试", "title2", "title3"}
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

	err = w.ReadUTF16Line(func(record []string) error {
		golog.WithTag("record").Info(record)
		return nil
	})
	if err != nil {
		golog.Error(err)
		return
	}
}

func TestCsvRead1(t *testing.T) {
	w, err := NewCsvReader("/Users/Jerry/Documents/my/test/meta/meta.csv", ',')
	if err != nil {
		golog.Error(err)
		return
	}
	line := 0
	err = w.ReadLineJson(UTF8, func(record map[string]string) error {
		golog.WithTag("record").Info(record)
		line++
		return nil
	})
	if err != nil {
		golog.Error(err)
		return
	}
	golog.Info(line)
}
