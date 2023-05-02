package services

import (
	"fmt"

	"mpj/internal/interfaces"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ interfaces.ILoggerService = (*LoggerService)(nil)

type LoggerService struct {
	// Name of the logger
	ZapLogger *zap.Logger
}

var logger *zap.Logger

func init() {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.FullCallerEncoder,
		},
	}
	zlogger, _ := cfg.Build(zap.AddCallerSkip(1))

	logger = zlogger

}

func NewLoggerService(typ, name string) *LoggerService {
	return &LoggerService{
		ZapLogger: logger.With(zap.String(typ, name)),
	}
}

func (svc *LoggerService) Logger() *zap.Logger {
	return svc.ZapLogger
}

func (svc *LoggerService) LogErr(err error) error {
	errMsg := fmt.Sprintf("error: %s", err.Error())
	svc.ZapLogger.Info(errMsg)

	return err
}

func (svc *LoggerService) Logf(level zapcore.Level, text string, format ...interface{}) {
	switch level {
	case zapcore.DebugLevel:
		svc.ZapLogger.Debug(fmt.Sprintf(text, format...))
	case zapcore.WarnLevel:
		svc.ZapLogger.Warn(fmt.Sprintf(text, format...))
	case zapcore.PanicLevel:
		svc.ZapLogger.Panic(fmt.Sprintf(text, format...))
	case zapcore.InfoLevel:
		svc.ZapLogger.Info(fmt.Sprintf(text, format...))
	case zapcore.FatalLevel:
		svc.ZapLogger.Fatal(fmt.Sprintf(text, format...))
	case zapcore.ErrorLevel:
		svc.ZapLogger.Error(fmt.Sprintf(text, format...))
	}
}
