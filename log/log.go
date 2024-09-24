package log

import (
	"os"
	"pcgamedb/config"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var ConsoleLogger *zap.Logger
var FileLogger *zap.Logger
var TaskLogger *zap.Logger

func init() {
	fileCore, consoleCore, combinedCore, taskCore := buildZapCore(getZapLogLevel(config.Config.LogLevel))
	FileLogger = zap.New(fileCore, zap.AddCaller())
	ConsoleLogger = zap.New(consoleCore, zap.AddCaller())
	Logger = zap.New(combinedCore, zap.AddCaller())
	TaskLogger = zap.New(taskCore, zap.AddCaller())
}

func buildZapCore(logLevel zapcore.Level) (fileCore zapcore.Core, consoleCore zapcore.Core, combinedCore zapcore.Core, taskCore zapcore.Core) {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})
	taskFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/task.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})
	consoleWriter := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	fileCore = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), fileWriter, logLevel)
	consoleCore = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleWriter, logLevel)
	combinedCore = zapcore.NewTee(fileCore, consoleCore)
	taskCore = zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), taskFileWriter, logLevel)
	return
}

func getZapLogLevel(logLevel string) zapcore.Level {
	switch strings.ToLower(logLevel) {
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "info":
		return zap.InfoLevel
	default:
		return zap.InfoLevel
	}
}
