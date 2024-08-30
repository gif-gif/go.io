package eslog

import (
	golog "github.com/gif-gif/go.io/go-log"
	"strings"
)

type Logger struct {
	Level golog.Level
}

func (l Logger) Printf(format string, v ...interface{}) {
	log := golog.WithTag("goo-es")
	switch l.Level {
	case golog.DEBUG:
		log.DebugF(format, v...)
	case golog.INFO:
		log.InfoF(format, v...)
	case golog.WARN, golog.ERROR, golog.PANIC, golog.FATAL:
		if strings.Contains(format, "warning") {
			log.WarnF(format, v...)
			return
		}
		log.ErrorF(format, v...)
	}
}
