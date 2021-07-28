package xlog

import (
	"sync"

	"github.com/op/go-logging"
)

// 日志级别
const (
	CRITICAL = logging.CRITICAL
	ERROR    = logging.ERROR
	WARNING  = logging.WARNING
	NOTICE   = logging.NOTICE
	INFO     = logging.INFO
	DEBUG    = logging.DEBUG
)

var (
	// 默认日志模块名称
	DefaultLoggerName string

	// 日志实例组
	loggers = make(map[string]*logger)

	// 全局监听器组
	listeners = make(map[logging.Level][]ListenerFunc)

	mu sync.Mutex
)

//------------------------------------------------------------------------------

type ListenerFunc func(module string, level logging.Level, format *string, args ...interface{})

// 注册全局监听器，作用于所有logger
func Listen(level logging.Level, fn ...ListenerFunc) {
	if _, ok := listeners[level]; !ok {
		listeners[level] = make([]ListenerFunc, 0)
	}
	listeners[level] = append(listeners[level], fn...)
}

//------------------------------------------------------------------------------

type logger struct {
	logger    *logging.Logger
	listeners map[logging.Level][]ListenerFunc
}

func (l *logger) Module() string {
	return l.logger.Module
}

func (l *logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *logger) Critical(args ...interface{}) {
	l.logger.Critical(args...)
	l.dispatch(CRITICAL, nil, args...)
}

func (l *logger) Criticalf(format string, args ...interface{}) {
	l.logger.Criticalf(format, args...)
	l.dispatch(CRITICAL, &format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.logger.Error(args...)
	l.dispatch(ERROR, nil, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
	l.dispatch(ERROR, &format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.logger.Warning(args...)
	l.dispatch(WARNING, nil, args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, args...)
	l.dispatch(WARNING, &format, args...)
}

func (l *logger) Notice(args ...interface{}) {
	l.logger.Notice(args...)
	l.dispatch(NOTICE, nil, args...)
}

func (l *logger) Noticef(format string, args ...interface{}) {
	l.logger.Noticef(format, args...)
	l.dispatch(NOTICE, &format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.logger.Info(args...)
	l.dispatch(INFO, nil, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
	l.dispatch(INFO, &format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
	l.dispatch(DEBUG, nil, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
	l.dispatch(DEBUG, &format, args...)
}

func (l *logger) Listen(level logging.Level, fn ...ListenerFunc) {
	if _, ok := l.listeners[level]; !ok {
		l.listeners[level] = make([]ListenerFunc, 0)
	}
	l.listeners[level] = append(l.listeners[level], fn...)
}

func (l *logger) dispatch(level logging.Level, format *string, args ...interface{}) {
	// 触发绑定的监听器
	if _, ok := l.listeners[level]; ok {
		for _, listener := range l.listeners[level] {
			listener(l.Module(), level, format, args...)
		}
	}

	// 触发全局的监听器
	if _, ok := listeners[level]; ok {
		for _, listener := range listeners[level] {
			listener(l.Module(), level, format, args...)
		}
	}
}

//------------------------------------------------------------------------------

func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

func Critical(args ...interface{}) {
	GetLogger().Critical(args...)
}

func Criticalf(format string, args ...interface{}) {
	GetLogger().Criticalf(format, args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Warning(args ...interface{}) {
	GetLogger().Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	GetLogger().Warningf(format, args...)
}

func Notice(args ...interface{}) {
	GetLogger().Notice(args...)
}

func Noticef(format string, args ...interface{}) {
	GetLogger().Noticef(format, args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

//------------------------------------------------------------------------------

func resolveLoggerName(names ...string) string {
	if len(names) > 0 && len(names[0]) > 0 {
		return names[0]
	}
	return DefaultLoggerName
}

func GetLogger(names ...string) *logger {
	name := resolveLoggerName(names...)
	if _, ok := loggers[name]; ok {
		return loggers[name]
	}
	mu.Lock()
	defer mu.Unlock()
	loggers[name] = &logger{
		logger:    logging.MustGetLogger(name),
		listeners: make(map[logging.Level][]ListenerFunc),
	}
	return loggers[name]
}
