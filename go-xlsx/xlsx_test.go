package goxlsx

import (
	golog "github.com/gif-gif/go.io/go-log"
	"testing"
)

func TestXlsxWrite(t *testing.T) {
	w := New()
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
