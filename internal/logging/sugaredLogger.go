package logging

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type sugaredLogger struct {
	devLogger *zap.Logger
	logger    *zap.SugaredLogger
}

func MakeSugaredLogger() Logger {
	logger := &sugaredLogger{}
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger.devLogger, _ = config.Build()
	logger.logger = logger.devLogger.Sugar()
	return logger
}

func (l sugaredLogger) Sync() {
	err := l.devLogger.Sync()
	if err != nil {
		log.Panic(err)
	}
}

func (l sugaredLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args)
}

func (l sugaredLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args)
}

func (l sugaredLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args)
}

func (l sugaredLogger) Warn(args ...interface{}) {
	l.logger.Warn(args)
}

func (l sugaredLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args)
}

func (l sugaredLogger) Debug(args ...interface{}) {
	l.logger.Debug(args)
}

func (l sugaredLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args)
}

func (l sugaredLogger) Error(args ...interface{}) {
	l.logger.Error(args)
}

func (l sugaredLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args)
}

func (l sugaredLogger) Info(args ...interface{}) {
	l.logger.Info(args)
}

func (l sugaredLogger) Panic(args ...interface{}) {
	l.logger.Panic(args)
}
