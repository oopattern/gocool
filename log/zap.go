package log

import (
	"os"
	"fmt"
	"runtime"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ZapLogger *zap.Logger
)

func init() {
	// set time, console, caller format
	cfg := zap.NewProductionEncoderConfig()
	// 文件名和行号
	cfg.CallerKey = "caller"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.LineEnding = zapcore.DefaultLineEnding
	encoder := zapcore.NewConsoleEncoder(cfg)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	// 开启开发模式, 堆栈跟踪
	caller := zap.AddCaller()
	// 文件名和行号
	callerSkip := zap.AddCallerSkip(1)
	development := zap.Development()
	ZapLogger = zap.New(core, caller, callerSkip, development)
	// 设置程序退出前强制刷盘
	runtime.SetFinalizer(ZapLogger, final)
}

// 强制刷磁盘
func final(logHandler *zap.Logger) {
	var _ = logHandler.Sync()
}

func Debug(format string, v ...interface{}) {
	ZapLogger.Debug(fmt.Sprintf(format, v))
}

func Info(format string, v ...interface{}) {
	ZapLogger.Info(fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	ZapLogger.Warn(fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	ZapLogger.Error(fmt.Sprintf(format, v...))
}

func Panic(format string, v ...interface{}) {
	ZapLogger.Panic(fmt.Sprintf(format, v...))
}

func Fatal(format string, v ...interface{}) {
	ZapLogger.Fatal(fmt.Sprintf(format, v...))
}


