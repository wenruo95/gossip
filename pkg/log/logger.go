package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string)
	Debugf(format string, v ...interface{})
	Debugw(format string, v ...zap.Field)
	Info(msg string)
	Infof(format string, v ...interface{})
	Infow(format string, v ...zap.Field)
	Warn(msg string)
	Warnf(format string, v ...interface{})
	Warnw(format string, v ...zap.Field)
	Error(msg string)
	Errorf(format string, v ...interface{})
	Errorw(format string, v ...zap.Field)
	Fatal(msg string)
	Fatalf(format string, v ...interface{})
	Fatalw(format string, v ...zap.Field)
	Level() zapcore.Level
}

type ZLogger struct {
	l *zap.Logger
}

func NewZLogger(l *zap.Logger) *ZLogger {
	return &ZLogger{l: l}
}

func (z *ZLogger) Debug(msg string) {
	z.l.Debug(msg)
}

func (z *ZLogger) Debugf(format string, v ...interface{}) {
	z.l.Debug(fmt.Sprintf(format, v...))
}

func (z *ZLogger) Debugw(format string, v ...zap.Field) {
	z.l.Debug(format, v...)
}

func (z *ZLogger) Info(msg string) {
	z.l.Info(msg)
}

func (z *ZLogger) Infof(format string, v ...interface{}) {
	z.l.Info(fmt.Sprintf(format, v...))
}

func (z *ZLogger) Infow(format string, v ...zap.Field) {
	z.l.Info(format, v...)
}

func (z *ZLogger) Warn(msg string) {
	z.l.Warn(msg)
}

func (z *ZLogger) Warnf(format string, v ...interface{}) {
	z.l.Warn(fmt.Sprintf(format, v...))
}

func (z *ZLogger) Warnw(format string, v ...zap.Field) {
	z.l.Warn(format, v...)
}

func (z *ZLogger) Error(msg string) {
	z.l.Error(msg)
}

func (z *ZLogger) Errorf(format string, v ...interface{}) {
	z.l.Error(fmt.Sprintf(format, v...))
}

func (z *ZLogger) Errorw(format string, v ...zap.Field) {
	z.l.Error(format, v...)
}

func (z *ZLogger) Fatal(msg string) {
	z.l.Fatal(msg)
}

func (z *ZLogger) Fatalf(format string, v ...interface{}) {
	z.l.Fatal(fmt.Sprintf(format, v...))
}

func (z *ZLogger) Fatalw(format string, v ...zap.Field) {
	z.l.Fatal(format, v...)
}

func (z *ZLogger) Level() zapcore.Level {
	return z.l.Level()
}
