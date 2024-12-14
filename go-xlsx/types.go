package goxlsx

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
)

// 常用文件编码
// unicode.UTF8,
// unicode.UTF16(BigEndian, UseBOM),
// unicode.UTF16(BigEndian, IgnoreBOM),
// unicode.UTF16(LittleEndian, IgnoreBOM),
// 常用文件编码需要用
var (
	UTF8    = unicode.UTF8
	UTF8BOM = unicode.UTF8BOM
	GBK     = simplifiedchinese.GBK
	//UTF16 有很多种 参考 unicode包
	UTF16 = unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
)

// 常用文件编码写入bom头部
var (
	BOM_UTF8  = []byte{0xEF, 0xBB, 0xBF}
	BOM_UTF16 = []byte{0xFF, 0xFE}
)
