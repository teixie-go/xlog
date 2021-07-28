package xlog

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	// 缓冲分隔符，默认为空格
	BufferSeparator = " "
)

type bufferLogger struct {
	name   string
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
	GetLogger(b.name).Info(b.String())
}

func appendString(args ...interface{}) string {
	str := ""
	for v := range args {
		str += fmt.Sprint(v) + BufferSeparator
	}
	return str
}

func NewDefaultBufferLogger(args ...interface{}) *bufferLogger {
	return &bufferLogger{buffer: bytes.NewBufferString(appendString(args...))}
}

func NewBufferLogger(name string, args ...interface{}) *bufferLogger {
	return &bufferLogger{
		name:   name,
		buffer: bytes.NewBufferString(appendString(args...)),
	}
}
