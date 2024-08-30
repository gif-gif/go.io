package golog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type ConsoleAdapter struct {
}

func NewConsoleLog() *Logger {
	return New(&ConsoleAdapter{})
}

func (ca *ConsoleAdapter) Write(msg *Message) {
	var (
		buf   bytes.Buffer
		nw    = time.Now()
		level = LevelText[msg.Level]
		color = Colors[msg.Level]
	)

	buf.WriteString(nw.Format("2006-01-02 15:04:05"))
	buf.WriteString(" ")

	buf.WriteString(color(fmt.Sprintf("%-5s", level)))
	buf.WriteString(" ")

	if l := len(msg.Entry.Tags); l > 0 {
		buf.WriteString("[" + strings.Join(msg.Entry.Tags, "][") + "]")
		buf.WriteString(" ")
	}

	if l := len(msg.Message); l > 0 {
		for _, v := range msg.Message {
			buf.WriteString(fmt.Sprint(v))
			buf.WriteString(" ")
		}
	}

	if l := len(msg.Entry.Data); l > 0 {
		data := map[string]interface{}{}
		for _, i := range msg.Entry.Data {
			data[i.Field] = i.Value
		}
		if b, err := json.Marshal(&data); err == nil {
			buf.Write(b)
		}
		buf.WriteString(" ")
	}

	if l := len(msg.Trace); l > 0 {
		if b, err := json.Marshal(&msg.Trace); err == nil {
			buf.Write(b)
		}
		buf.WriteString(" ")
	}

	ca.writer().Write(append(buf.Bytes(), '\n'))
}

func (ca ConsoleAdapter) writer() io.Writer {
	return os.Stdout
}
