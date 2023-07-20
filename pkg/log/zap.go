package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	zapCore   zapcore.Core
	stdLogger Logger
)

func init() {
	zapCore = zapcore.NewCore(getEncoder(), getSyncer(nil), zapcore.DebugLevel)
	stdLogger = NewZLogger(ZapLogger(zap.AddCallerSkip(2)))
}

func ZapCore() zapcore.Core {
	return zapCore
}

func ZapLogger(options ...zap.Option) *zap.Logger {
	options = append(options, zap.AddCaller())
	return zap.New(zapCore, options...)
}

func CtxLogger() *ZLogger {
	return NewZLogger(ZapLogger(zap.AddCallerSkip(2)))
}

type Config struct {
	FileName   string // log output path
	Level      string // debug info warn error fatal, default debug
	MaxSize    int    // log maxsize before rotated, default 100M
	MaxAge     int    // reserve log max time day, default: 30 day
	MaxBackups int    // reserve old log file number, default 100
	LocalTime  bool   // default use UTC TIME
	Compress   bool   // use gzip to compress
}

func InitLogger(config *Config) error {
	return initZapLog(config)
}

func Sync() {
	if err := zapCore.Sync(); err != nil {
		fmt.Printf("sync error:" + err.Error())
	}
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getSyncer(config *Config) zapcore.WriteSyncer {
	if config == nil || len(config.FileName) == 0 {
		return zapcore.AddSync(os.Stdout)
	}

	writer := &lumberjack.Logger{
		Filename:   config.FileName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	return zapcore.AddSync(writer)
}

func initZapLog(config *Config) error {
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return err
	}

	zapCore = zapcore.NewCore(getEncoder(), getSyncer(config), level)
	stdLogger = NewZLogger(ZapLogger(zap.AddCallerSkip(2)))
	return nil
}
