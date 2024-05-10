package golog

import "sync"

var (
	__log  *Logger
	__once sync.Once
)

func Default() *Logger {
	__once.Do(func() {
		__log = NewConsoleLog()
	})
	return __log
}

func SetAdapter(adapter Adapter) {
	Default().adapter = adapter
}

func WithHook(fns ...func(msg *Message)) {
	Default().WithHook(fns...)
}

func WithTag(tags ...string) *Entry {
	return Default().WithTag(tags...)
}

func WithField(field string, value interface{}) *Entry {
	return Default().WithField(field, value)
}

func Debug(v ...interface{}) {
	Default().Debug(v...)
}

func DebugF(format string, v ...interface{}) {
	Default().DebugF(format, v...)
}

func Info(v ...interface{}) {
	Default().Info(v...)
}

func InfoF(format string, v ...interface{}) {
	Default().InfoF(format, v...)
}

func Warn(v ...interface{}) {
	Default().Warn(v...)
}

func WarnF(format string, v ...interface{}) {
	Default().WarnF(format, v...)
}

func Error(v ...interface{}) {
	Default().Error(v...)
}

func ErrorF(format string, v ...interface{}) {
	Default().ErrorF(format, v...)
}

func Panic(v ...interface{}) {
	Default().Panic(v...)
}

func PanicF(format string, v ...interface{}) {
	Default().PanicF(format, v...)
}

func Fatal(v ...interface{}) {
	Default().Fatal(v...)
}

func FatalF(format string, v ...interface{}) {
	Default().FatalF(format, v...)
}
