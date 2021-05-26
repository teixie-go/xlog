package xlog

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	BufferSeparator = " "
)

type bufferLogger struct {
	buffer *bytes.Buffer
}

func (b *bufferLogger) Append(args ...interface{}) {
	b.buffer.WriteString(appendString(args...))
}

func (b *bufferLogger) Appendf(format string, args ...interface{}) {
	b.Append(fmt.Sprintf(format, args...))
}

func (b *bufferLogger) String() string {
	return strings.TrimRight(b.buffer.String(), BufferSeparator)
}

func (b *bufferLogger) Flush() {
	Info(b.String())
}

func appendString(args ...interface{}) string {
	str := ""
	for v := range args {
		str += fmt.Sprint(v) + BufferSeparator
	}
	return str
}

func NewBufferLogger(args ...interface{}) *bufferLogger {
	return &bufferLogger{buffer: bytes.NewBufferString(appendString(args...))}
}
