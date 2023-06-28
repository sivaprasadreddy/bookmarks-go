package logging

import (
	"github.com/sivaprasadreddy/bookmarks-go/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(cfg config.AppConfig) *Logger {
	return initZap(cfg)
}

func initZap(cfg config.AppConfig) *Logger {
	logFile := "bookmarks.log"
	logLevel := zap.DebugLevel
	hook := lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1024,
		MaxBackups: 30,
		MaxAge:     7,
		Compress:   true,
	}

	encoder := getEncoder()
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		logLevel)
	if cfg.Environment != "prod" {
		return &Logger{zap.New(
			core, zap.Development(),
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
		).Sugar()}
	}
	return &Logger{zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	).Sugar()}
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:   "ts",
		LevelKey:  "level",
		NameKey:   "logger",
		CallerKey: "caller",
		//FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}
