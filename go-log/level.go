package golog

type Level int
type brush func(string) string

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

var (
	LevelText = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		PANIC: "PANIC",
		FATAL: "FATAL",
	}

	Colors = map[Level]brush{
		DEBUG: newBrush("1;34"), // blue
		INFO:  newBrush("1;32"), // green
		WARN:  newBrush("1;33"), // yellow
		ERROR: newBrush("1;31"), // red
		PANIC: newBrush("1;37"), // white
		FATAL: newBrush("1;35"), // magenta
	}
)

func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}
