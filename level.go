package xlog

// Level levels
type Level int8

const (
	DEBUG Level = iota - 1
	INFO
	WARNING
	ERROR
	PANIC
	FATAL
)

var levelNames = map[Level]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	PANIC:   "PANIC",
	FATAL:   "FATAL",
}

func (l Level) String() string {
	return levelNames[l]
}
