package xlog

import (
	"context"

	"go.uber.org/zap"
)

var (
	_ Handler = (*zapHandler)(nil)
)

type zapHandler struct {
	logger *zap.Logger
}

func (h *zapHandler) Log(ctx context.Context, params Params) {
	l := h.logger.Sugar()
	if len(params.Fields) > 0 {
		l = l.With(params.Fields...)
	}
	template := ""
	if params.Format != nil {
		template = *params.Format
	}
	switch params.Level {
	case DEBUG:
		l.Debugf(template, params.Args...)
	case INFO:
		l.Infof(template, params.Args...)
	case WARNING:
		l.Warnf(template, params.Args...)
	case ERROR:
		l.Errorf(template, params.Args...)
	case PANIC:
		l.Panicf(template, params.Args...)
	case FATAL:
		l.Fatalf(template, params.Args...)
	default:
	}
}

func NewZapHandler(logger *zap.Logger) Handler {
	return &zapHandler{logger: logger}
}
