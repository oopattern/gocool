package log

import (
	"os"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ZapLogger *zap.Logger
)

func init() {
	// set time, console, caller format
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	ZapLogger = zap.New(core, zap.AddCaller())
}

func Debug(format string, v ...interface{}) {
	ZapLogger.Debug(fmt.Sprintf(format, v))
}

func Info(format string, v ...interface{}) {
	ZapLogger.Info(fmt.Sprintf(format, v))
}

func Warn(format string, v ...interface{}) {
	ZapLogger.Warn(fmt.Sprintf(format, v))
}

func Error(format string, v ...interface{}) {
	ZapLogger.Error(fmt.Sprintf(format, v))
}

func Panic(format string, v ...interface{}) {
	ZapLogger.Panic(fmt.Sprintf(format, v))
}

func Fatal(format string, v ...interface{}) {
	ZapLogger.Fatal(fmt.Sprintf(format, v))
}


