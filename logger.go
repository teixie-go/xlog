package xlog

import (
	"context"
	"path/filepath"
	"runtime"
)

const (
	CallerSkipOffset = 3
)

var (
	// global middleware
	_middleware Middleware
)

type (
	logger struct {
		addCaller  bool
		callerSkip int
		handlers   []Handler
		middleware Middleware
	}

	Handler interface {
		Log(ctx context.Context, params Params)
	}

	Params struct {
		Caller *Caller
		Level  Level
		Format *string
		Args   []interface{}
		Fields []interface{}
	}

	Param func(*Params)

	Option func(*logger)
)

func (l *logger) Fatal(ctx context.Context, args ...interface{}) {
	l.log(ctx, FATAL, nil, Args(args...))
}

func (l *logger) Fatalv(ctx context.Context, param ...Param) {
	l.log(ctx, FATAL, nil, param...)
}

func (l *logger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, FATAL, &format, Args(args...))
}

func (l *logger) Panic(ctx context.Context, args ...interface{}) {
	l.log(ctx, PANIC, nil, Args(args...))
}

func (l *logger) Panicv(ctx context.Context, param ...Param) {
	l.log(ctx, PANIC, nil, param...)
}

func (l *logger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, PANIC, &format, Args(args...))
}

func (l *logger) Error(ctx context.Context, args ...interface{}) {
	l.log(ctx, ERROR, nil, Args(args...))
}

func (l *logger) Errorv(ctx context.Context, param ...Param) {
	l.log(ctx, ERROR, nil, param...)
}

func (l *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, ERROR, &format, Args(args...))
}

func (l *logger) Warning(ctx context.Context, args ...interface{}) {
	l.log(ctx, WARNING, nil, Args(args...))
}

func (l *logger) Warningv(ctx context.Context, param ...Param) {
	l.log(ctx, WARNING, nil, param...)
}

func (l *logger) Warningf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, WARNING, &format, Args(args...))
}

func (l *logger) Info(ctx context.Context, args ...interface{}) {
	l.log(ctx, INFO, nil, Args(args...))
}

func (l *logger) Infov(ctx context.Context, param ...Param) {
	l.log(ctx, INFO, nil, param...)
}

func (l *logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, INFO, &format, Args(args...))
}

func (l *logger) Debug(ctx context.Context, args ...interface{}) {
	l.log(ctx, DEBUG, nil, Args(args...))
}

func (l *logger) Debugv(ctx context.Context, param ...Param) {
	l.log(ctx, DEBUG, nil, param...)
}

func (l *logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, DEBUG, &format, Args(args...))
}

func (l *logger) log(ctx context.Context, level Level, format *string, param ...Param) {
	params := Params{
		Level:  level,
		Format: format,
		Args:   make([]interface{}, 0),
		Fields: make([]interface{}, 0),
	}
	for _, p := range param {
		p(&params)
	}
	if l.addCaller {
		params.Caller = GetCaller(l.callerSkip + CallerSkipOffset)
	}
	if l.middleware != nil || _middleware != nil {
		closure := func(ctx context.Context, params *Params) {}
		if l.middleware != nil {
			closure = l.middleware(closure)
		}
		if _middleware != nil {
			closure = _middleware(closure)
		}
		closure(ctx, &params)
	}
	for _, h := range l.handlers {
		h.Log(ctx, params)
	}
}

func WithHandler(handlers ...Handler) Option {
	return func(l *logger) {
		l.handlers = append(l.handlers, handlers...)
	}
}

func WithMiddleware(middleware ...Middleware) Option {
	return func(l *logger) {
		l.middleware = withMiddlewareChain(l.middleware, middleware...)
	}
}

func WithCaller(enabled bool) Option {
	return func(l *logger) {
		l.addCaller = enabled
	}
}

func WithCallerSkip(skip int) Option {
	return func(l *logger) {
		l.callerSkip += skip
	}
}

func Args(args ...interface{}) Param {
	return func(p *Params) {
		p.Args = append(p.Args, args...)
	}
}

func Argsf(format string, args ...interface{}) Param {
	return func(p *Params) {
		if p.Format == nil {
			p.Format = &format
		} else {
			*p.Format += format
		}
		p.Args = append(p.Args, args...)
	}
}

func Fields(fields ...interface{}) Param {
	return func(p *Params) {
		p.Fields = append(p.Fields, fields...)
	}
}

func NewLogger(options ...Option) *logger {
	l := &logger{handlers: make([]Handler, 0)}
	for _, o := range options {
		o(l)
	}
	return l
}

//------------------------------------------------------------------------------

type (
	Closure func(ctx context.Context, params *Params)

	Middleware func(Closure) Closure
)

func withMiddlewareChain(chain Middleware, middleware ...Middleware) Middleware {
	if len(middleware) == 0 {
		return chain
	}
	if chain == nil {
		return withMiddlewareChain(middleware[0], middleware[1:]...)
	}
	return withMiddlewareChain(func(next Closure) Closure {
		return chain(middleware[0](next))
	}, middleware[1:]...)
}

func Use(middleware ...Middleware) {
	_middleware = withMiddlewareChain(_middleware, middleware...)
}

//------------------------------------------------------------------------------

type Caller struct {
	PC       uintptr
	File     string
	Filename string
	Function string
	Line     int
}

func GetCaller(skip int) *Caller {
	function := "???"
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}
	return &Caller{
		PC:       pc,
		File:     file,
		Filename: filepath.Base(file),
		Function: function,
		Line:     line,
	}
}
