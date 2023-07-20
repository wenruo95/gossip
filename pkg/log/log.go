package log

import (
	"context"
	"fmt"

	"github.com/wenruo95/gossip/pkg/codec"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Debug(msg string) {
	stdLogger.Debug(msg)
}

func Debugf(format string, v ...interface{}) {
	stdLogger.Debug(fmt.Sprintf(format, v...))
}

func Debugw(format string, v ...zap.Field) {
	stdLogger.Debugw(format, v...)
}

func Info(msg string) {
	stdLogger.Info(msg)
}

func Infof(format string, v ...interface{}) {
	stdLogger.Info(fmt.Sprintf(format, v...))
}

func Infow(format string, v ...zap.Field) {
	stdLogger.Infow(format, v...)
}

func Warn(msg string) {
	stdLogger.Warn(msg)
}

func Warnf(format string, v ...interface{}) {
	stdLogger.Warn(fmt.Sprintf(format, v...))
}

func Warnw(format string, v ...zap.Field) {
	stdLogger.Warnw(format, v...)
}

func Error(msg string) {
	stdLogger.Error(msg)
}

func Errorf(format string, v ...interface{}) {
	stdLogger.Error(fmt.Sprintf(format, v...))
}

func Errorw(format string, v ...zap.Field) {
	stdLogger.Errorw(format, v...)
}

func Fatal(msg string) {
	stdLogger.Fatal(msg)
}

func Fatalf(format string, v ...interface{}) {
	stdLogger.Fatal(fmt.Sprintf(format, v...))
}

func Fatalw(format string, v ...zap.Field) {
	stdLogger.Fatalw(format, v...)
}

func Level() zapcore.Level {
	return stdLogger.Level()
}

func CtxDebug(ctx context.Context, msg string) {
	stdLogger.Debugw(msg, codec.Message(ctx).Fileds()...)
}

func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	stdLogger.Debugw(fmt.Sprintf(format, v...), codec.Message(ctx).Fileds()...)
}

func CtxDebugw(ctx context.Context, format string, v ...zap.Field) {
	fs := make([]zap.Field, 0)
	fs = append(fs, codec.Message(ctx).Fileds()...)
	fs = append(fs, v...)
	stdLogger.Debugw(format, fs...)
}

func CtxInfo(ctx context.Context, msg string) {
	stdLogger.Infow(msg, codec.Message(ctx).Fileds()...)
}

func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	stdLogger.Infow(fmt.Sprintf(format, v...), codec.Message(ctx).Fileds()...)
}

func CtxInfow(ctx context.Context, format string, v ...zap.Field) {
	fs := make([]zap.Field, 0)
	fs = append(fs, codec.Message(ctx).Fileds()...)
	fs = append(fs, v...)
	stdLogger.Infow(format, fs...)
}

func CtxWarn(ctx context.Context, msg string) {
	stdLogger.Warnw(msg, codec.Message(ctx).Fileds()...)
}

func CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	stdLogger.Warnw(fmt.Sprintf(format, v...), codec.Message(ctx).Fileds()...)
}

func CtxWarnw(ctx context.Context, format string, v ...zap.Field) {
	fs := make([]zap.Field, 0)
	fs = append(fs, codec.Message(ctx).Fileds()...)
	fs = append(fs, v...)
	stdLogger.Warnw(format, fs...)
}

func CtxError(ctx context.Context, msg string) {
	stdLogger.Errorw(msg, codec.Message(ctx).Fileds()...)
}

func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	stdLogger.Errorw(fmt.Sprintf(format, v...), codec.Message(ctx).Fileds()...)
}

func CtxErrorw(ctx context.Context, format string, v ...zap.Field) {
	fs := make([]zap.Field, 0)
	fs = append(fs, codec.Message(ctx).Fileds()...)
	fs = append(fs, v...)
	stdLogger.Errorw(format, fs...)
}

func CtxFatal(ctx context.Context, msg string) {
	stdLogger.Fatalw(msg, codec.Message(ctx).Fileds()...)
}

func CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	stdLogger.Fatalw(fmt.Sprintf(format, v...), codec.Message(ctx).Fileds()...)
}

func CtxFatalw(ctx context.Context, format string, v ...zap.Field) {
	fs := make([]zap.Field, 0)
	fs = append(fs, codec.Message(ctx).Fileds()...)
	fs = append(fs, v...)
	stdLogger.Fatalw(format, fs...)
}
