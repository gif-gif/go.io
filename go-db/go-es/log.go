package goes

import (
	golog "github.com/gif-gif/go.io/go-log"
	"strings"
)

type logger struct {
	level golog.Level
}

func (l logger) Printf(format string, v ...interface{}) {
	log := golog.WithTag("goo-es")
	switch l.level {
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
