package goxlsx

import "os"

type CsvWriter struct {
	FilePath string
	Comma    rune //csv 列分割符
	file     *os.File
}
