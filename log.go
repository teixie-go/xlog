package xlog

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = &logger{}

type Config struct {
	Path       string `yaml:"path" json:"path"`
	Level      string `yaml:"level" json:"level"`
	MaxSize    int    `yaml:"max_size" json:"max_size"`
	MaxAge     int    `yaml:"max_age" json:"max_age"`
	MaxBackups int    `yaml:"max_backups" json:"max_backups"`
}

func Init(options ...Option) error {
	log = NewLogger(options...)
	return nil
}

func InitDefault(cfg Config) error {
	ls := zapcore.AddSync(os.Stdout)
	if len(cfg.Path) != 0 {
		ls = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Path,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
		})
	}

	// default level INFO
	level := zapcore.InfoLevel
	zapLevel := map[string]zapcore.Level{
		"DEBUG":   zapcore.DebugLevel,
		"INFO":    zapcore.InfoLevel,
		"WARNING": zapcore.WarnLevel,
		"ERROR":   zapcore.ErrorLevel,
		"PANIC":   zapcore.PanicLevel,
		"FATAL":   zapcore.FatalLevel,
	}
	if lvl, ok := zapLevel[strings.ToUpper(cfg.Level)]; ok {
		level = lvl
	}

	conf := zap.NewProductionEncoderConfig()
	conf.TimeKey = "t"
	conf.LevelKey = "l"
	conf.CallerKey = "c"
	conf.EncodeLevel = capitalLevelEncoder
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	enc := zapcore.NewJSONEncoder(conf)
	zapLogger := zap.New(zapcore.NewCore(enc, ls, level), zap.AddCaller(), zap.AddCallerSkip(CallerSkipOffset+1))
	zap.ReplaceGlobals(zapLogger)

	return Init(WithHandler(NewZapHandler(zapLogger)), WithCaller(true), WithCallerSkip(1))
}

func capitalLevelEncoder(lvl zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch lvl {
	case zapcore.DebugLevel:
		enc.AppendString("D")
	case zapcore.InfoLevel:
		enc.AppendString("I")
	case zapcore.WarnLevel:
		enc.AppendString("W")
	case zapcore.ErrorLevel:
		enc.AppendString("E")
	case zapcore.DPanicLevel:
		enc.AppendString("DP")
	case zapcore.PanicLevel:
		enc.AppendString("P")
	case zapcore.FatalLevel:
		enc.AppendString("F")
	default:
		enc.AppendString(fmt.Sprintf("LEVEL(%d)", lvl))
	}
}

func Fatal(ctx context.Context, args ...interface{}) {
	log.Fatal(ctx, args...)
}

func Fatalv(ctx context.Context, param ...Param) {
	log.Fatalv(ctx, param...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	log.Fatalf(ctx, format, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	log.Panic(ctx, args...)
}

func Panicv(ctx context.Context, param ...Param) {
	log.Panicv(ctx, param...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	log.Panicf(ctx, format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	log.Error(ctx, args...)
}

func Errorv(ctx context.Context, param ...Param) {
	log.Errorv(ctx, param...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	log.Warning(ctx, args...)
}

func Warningv(ctx context.Context, param ...Param) {
	log.Warningv(ctx, param...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(ctx, format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	log.Info(ctx, args...)
}

func Infov(ctx context.Context, param ...Param) {
	log.Infov(ctx, param...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	log.Debug(ctx, args...)
}

func Debugv(ctx context.Context, param ...Param) {
	log.Debugv(ctx, param...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	log.Debugf(ctx, format, args...)
}
