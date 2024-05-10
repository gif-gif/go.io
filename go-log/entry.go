package golog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Entry struct {
	Tags []string
	Data []DataField
	msg  *Message
	l    *Logger
}

type DataField struct {
	Field string
	Value interface{}
}

func NewEntry(l *Logger) *Entry {
	return &Entry{l: l}
}

func (entry *Entry) WithTag(tags ...string) *Entry {
	entry.Tags = append(entry.Tags, tags...)
	return entry
}

func (entry *Entry) WithField(field string, value interface{}) *Entry {
	entry.Data = append(entry.Data, DataField{Field: field, Value: value})
	return entry
}

func (entry *Entry) Debug(v ...interface{}) {
	entry.output(DEBUG, v...)
}

func (entry *Entry) DebugF(format string, v ...interface{}) {
	entry.output(DEBUG, fmt.Sprintf(format, v...))
}

func (entry *Entry) Info(v ...interface{}) {
	entry.output(INFO, v...)
}

func (entry *Entry) InfoF(format string, v ...interface{}) {
	entry.output(INFO, fmt.Sprintf(format, v...))
}

func (entry *Entry) Warn(v ...interface{}) {
	entry.output(WARN, v...)
}

func (entry *Entry) WarnF(format string, v ...interface{}) {
	entry.output(WARN, fmt.Sprintf(format, v...))
}

func (entry *Entry) Error(v ...interface{}) {
	entry.output(ERROR, v...)
}

func (entry *Entry) ErrorF(format string, v ...interface{}) {
	entry.output(ERROR, fmt.Sprintf(format, v...))
}

func (entry *Entry) Panic(v ...interface{}) {
	entry.output(PANIC, v...)
}

func (entry *Entry) PanicF(format string, v ...interface{}) {
	entry.output(PANIC, fmt.Sprintf(format, v...))
}

func (entry *Entry) Fatal(v ...interface{}) {
	entry.output(FATAL, v...)
	os.Exit(1)
}

func (entry *Entry) FatalF(format string, v ...interface{}) {
	entry.output(FATAL, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (entry *Entry) output(level Level, v ...interface{}) {
	entry.msg = &Message{
		Level:   level,
		Message: v,
		Time:    time.Now(),
		Entry:   entry,
	}

	if level >= WARN {
		entry.msg.Trace = entry.trace()
	}

	for _, fn := range entry.l.hooks {
		go entry.hookHandler(fn)
	}

	if entry.l.adapter != nil {
		entry.l.adapter.Write(entry.msg)
	}
}

func (entry *Entry) hookHandler(fn func(msg *Message)) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	fn(entry.msg)
}

// runtime.Caller 仅能获取非 goroutine 的信息
func (entry *Entry) trace() (arr []string) {
	arr = []string{}

	for i := 3; i < 16; i++ {
		_, file, line, _ := runtime.Caller(i)
		if file == "" {
			continue
		}
		if strings.Contains(file, ".pb.go") ||
			strings.Contains(file, "runtime/") ||
			strings.Contains(file, "src/") ||
			strings.Contains(file, "pkg/mod/") ||
			strings.Contains(file, "vendor/") ||
			strings.Contains(file, "go-log") {
			continue
		}
		arr = append(arr, fmt.Sprintf("%s %dL", entry.prettyFile(file), line))
	}

	return
}

func (entry *Entry) prettyFile(file string) string {
	var (
		index  int
		index2 int
	)

	if index = strings.LastIndex(file, "src/test/"); index >= 0 {
		return file[index+9:]
	}
	if index = strings.LastIndex(file, "src/"); index >= 0 {
		return file[index+4:]
	}
	if index = strings.LastIndex(file, "pkg/mod/"); index >= 0 {
		return file[index+8:]
	}
	if index = strings.LastIndex(file, "vendor/"); index >= 0 {
		return file[index+7:]
	}

	if index = strings.LastIndex(file, "/"); index < 0 {
		return file
	}
	if index2 = strings.LastIndex(file[:index], "/"); index2 < 0 {
		return file[index+1:]
	}
	return file[index2+1:]
}
