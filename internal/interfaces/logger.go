package interfaces

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILoggerService interface {
	Logf(level zapcore.Level, text string, format ...interface{})
	LogErr(err error) error

	Logger() *zap.Logger
}
