package log

import (
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"chat/pkg/log/conf"
)

var _ log.Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	logger *zap.Logger
	Sync   func() error
}

func NewZapLogger(mode string, filePath string, maxSize, maxBackups, maxAge int, compress bool) *ZapLogger {
	// 打印错误级别的日志
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	// 打印所有级别的日志
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	var allCore []zapcore.Core

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stack"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if mode == "dev" {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		allCore = append(allCore, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), lowPriority))
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		//info文件writeSyncer
		infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filePath + "info.log", //日志文件存放目录，如果文件夹不存在会自动创建
			MaxSize:    maxSize,               //文件大小限制,单位MB
			MaxBackups: maxBackups,            //最大保留日志文件数量
			MaxAge:     maxAge,                //日志文件保留天数
			Compress:   compress,              //是否压缩处理
		})
		//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
		//infoFileCore := zapcore.NewCore(consoleEncoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer,zapcore.AddSync(os.Stdout)), lowPriority)
		infoFileCore := zapcore.NewCore(consoleEncoder, infoFileWriteSyncer, lowPriority)

		//error文件writeSyncer
		errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filePath + "error.log", //日志文件存放目录
			MaxSize:    maxSize,                //文件大小限制,单位MB
			MaxBackups: maxBackups,             //最大保留日志文件数量
			MaxAge:     maxAge,                 //日志文件保留天数
			Compress:   compress,               //是否压缩处理
		})
		//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志
		errorFileCore := zapcore.NewCore(consoleEncoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

		allCore = append(allCore, infoFileCore)
		allCore = append(allCore, errorFileCore)
	}

	core := zapcore.NewTee(allCore...)

	Logger := zap.New(core).WithOptions(zap.AddCaller())
	zap.ReplaceGlobals(Logger)

	return &ZapLogger{logger: Logger, Sync: Logger.Sync}
}

// Log Implementation of Logger interface
func (l *ZapLogger) Log(level log.Level, keyVals ...interface{}) error {
	if len(keyVals) == 0 || len(keyVals)%2 != 0 {
		l.logger.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyVals))
		return nil
	}
	// Zap.Field is used when keyVals pairs appear
	var data []zap.Field
	for i := 0; i < len(keyVals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyVals[i]), fmt.Sprint(keyVals[i+1])))
	}
	switch level {
	case log.LevelDebug:
		l.logger.Debug("", data...)
	case log.LevelInfo:
		l.logger.Info("", data...)
	case log.LevelWarn:
		l.logger.Warn("", data...)
	case log.LevelError:
		l.logger.Error("", data...)
	}
	return nil
}

func (l *ZapLogger) GetZapLogger() *zap.Logger {
	return l.logger
}

func (l *ZapLogger) GetLogger(id, name, version string) log.Logger {
	logger := interface{}(l).(log.Logger)
	logger = log.With(logger, "trace_id", tracing.TraceID())
	logger = log.With(logger, "span_id", tracing.SpanID())
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "service.id", id)
	logger = log.With(logger, "service.name", name)
	logger = log.With(logger, "service.version", version)
	return logger
}

func Logger(mode string, zapLog *conf.Log_ZapLog) *ZapLogger {
	return NewZapLogger(mode, zapLog.LogPath, int(zapLog.MaxSize), int(zapLog.MaxBackups), int(zapLog.MaxAge), zapLog.Compress)
}
