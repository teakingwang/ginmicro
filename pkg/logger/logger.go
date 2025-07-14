package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var log *zap.Logger
var sugar *zap.SugaredLogger

// Init 初始化日志，development 控制是否用开发模式
func Init(development bool) error {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // INFO, ERROR 等大写彩色
		EncodeTime:     customTimeEncoder,                // 使用自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoder),       // 编码器
		zapcore.AddSync(os.Stdout),               // 输出
		zap.NewAtomicLevelAt(zapcore.DebugLevel), // 日志级别
	)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar = log.Sugar()
	return nil
}

// 自定义时间格式为 datetime（例如：2025-06-08 14:35:20）
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// 通用封装
func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Infof(msg string, args ...interface{}) {
	sugar.Infof(msg, args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Warnf(msg string, args ...interface{}) {
	sugar.Warnf(msg, args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func Errorf(msg string, args ...interface{}) {
	sugar.Errorf(msg, args...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Fatalf(msg string, args ...interface{}) {
	sugar.Fatalf(msg, args...)
}

func Debugf(msg string, args ...interface{}) {
	sugar.Debugf(msg, args...)
}

func Sync() {
	_ = log.Sync()
}
