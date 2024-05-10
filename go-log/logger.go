package golog

type Logger struct {
	hooks   []func(msg *Message)
	adapter Adapter
}

func New(adapter Adapter) *Logger {
	return &Logger{
		adapter: adapter,
	}
}

func (l *Logger) SetAdapter(adapter Adapter) {
	l.adapter = adapter
}

func (l *Logger) WithHook(fns ...func(msg *Message)) {
	l.hooks = append(l.hooks, fns...)
}

func (l *Logger) WithTag(tags ...string) *Entry {
	return NewEntry(l).WithTag(tags...)
}

func (l *Logger) WithField(field string, value interface{}) *Entry {
	return NewEntry(l).WithField(field, value)
}

func (l *Logger) Debug(v ...interface{}) {
	NewEntry(l).Debug(v...)
}

func (l *Logger) DebugF(format string, v ...interface{}) {
	NewEntry(l).DebugF(format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	NewEntry(l).Info(v...)
}

func (l *Logger) InfoF(format string, v ...interface{}) {
	NewEntry(l).InfoF(format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	NewEntry(l).Warn(v...)
}

func (l *Logger) WarnF(format string, v ...interface{}) {
	NewEntry(l).WarnF(format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	NewEntry(l).Error(v...)
}

func (l *Logger) ErrorF(format string, v ...interface{}) {
	NewEntry(l).ErrorF(format, v...)
}

func (l *Logger) Panic(v ...interface{}) {
	NewEntry(l).Panic(v...)
}

func (l *Logger) PanicF(format string, v ...interface{}) {
	NewEntry(l).PanicF(format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	NewEntry(l).Fatal(v...)
}

func (l *Logger) FatalF(format string, v ...interface{}) {
	NewEntry(l).FatalF(format, v...)
}
