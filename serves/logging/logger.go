package logging

import (
	"FlyCloud/pkg/system"
	"FlyCloud/serves/config"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

// 声明一个日志对象
var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

// 初始化日志对象
func InitLogger(cfg *config.LoggerConfig) {
	fmt.Println("------------init logger----------")
	var (
		err error
	)
	// 构造日志配置
	EncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	// 判断日志路径是否存在
	if _, err = system.IsExist(cfg.Path); err != nil {
		// 如果不存在，则创建日志目录
		if err = system.MkDir(cfg.Path); err != nil {
			fmt.Println("创建日志目录失败！", err)
			return
		}
	}
	// 构造Config
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     EncoderConfig,
		InitialFields:     map[string]interface{}{"MyName": "yueqing"},
		OutputPaths:       []string{"stdout", cfg.Path},
		ErrorOutputPaths:  []string{"stderr"},
	}

	// 构造Logger
	Logger, err = config.Build()
	if err != nil {
		log.Println("初始化日志失败！", err)
	}
	// 构造SugarLogger
	SugarLogger = Logger.Sugar()

	fmt.Println("------------init logger success----------")
}

// 失败日志
func Error(args ...interface{}) {
	SugarLogger.Error(args...)
}

// Debug日志
func Debug(args ...interface{}) {
	SugarLogger.Debug(args...)
}

// Info日志
func Info(args ...interface{}) {
	SugarLogger.Info(args...)
}

// Warn日志
func Warn(args ...interface{}) {
	SugarLogger.Warn(args...)
}

// Fatal日志
func Fatal(args ...interface{}) {
	SugarLogger.Fatal(args...)
}

// Panic日志
func Panic(args ...interface{}) {
	SugarLogger.Panic(args...)
}

// 记录日志
func Log(level string, args ...interface{}) {
	switch level {
	case "debug":
		SugarLogger.Debug(args...)
	case "info":
		SugarLogger.Info(args...)
	case "warn":
		SugarLogger.Warn(args...)
	case "error":
		SugarLogger.Error(args...)
	case "fatal":
		SugarLogger.Fatal(args...)
	case "panic":
		SugarLogger.Panic(args...)
	default:
		SugarLogger.Info(args...)
	}
}
