package xlog

import (
	"github.com/op/go-logging"
)

type handlerFunc func(format *string, args ...interface{})

var handlers = make(map[logging.Level][]handlerFunc)

func Critical(args ...interface{}) {
	logger.Critical(args...)
	handle(logging.CRITICAL, nil, args...)
}

func Criticalf(format string, args ...interface{}) {
	logger.Criticalf(format, args...)
	handle(logging.CRITICAL, &format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
	handle(logging.ERROR, nil, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
	handle(logging.ERROR, &format, args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
	handle(logging.WARNING, nil, args...)
}

func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
	handle(logging.WARNING, &format, args...)
}

func Notice(args ...interface{}) {
	logger.Notice(args...)
	handle(logging.NOTICE, nil, args...)
}

func Noticef(format string, args ...interface{}) {
	logger.Noticef(format, args...)
	handle(logging.NOTICE, &format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
	handle(logging.INFO, nil, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
	handle(logging.INFO, &format, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
	handle(logging.DEBUG, nil, args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
	handle(logging.DEBUG, &format, args...)
}

func handle(level logging.Level, format *string, args ...interface{}) {
	functions, ok := handlers[level]
	if !ok || len(functions) <= 0 {
		return
	}
	for _, fn := range functions {
		fn(format, args...)
	}
}

func RegisterHandler(level logging.Level, fn handlerFunc) {
	if functions, ok := handlers[level]; ok {
		functions = append(functions, fn)
	} else {
		handlers[level] = []handlerFunc{fn}
	}
}

func RegisterErrorHandler(fn handlerFunc) {
	RegisterHandler(logging.ERROR, fn)
}

func RegisterWarningHandler(fn handlerFunc) {
	RegisterHandler(logging.WARNING, fn)
}
