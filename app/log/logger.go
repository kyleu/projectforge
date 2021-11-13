package log

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const keyCustom = "custom"

func InitLogging(debug bool) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error
	switch {
	case strings.EqualFold(os.Getenv("LOGGING_FORMAT"), "json"):
		logger, err = initJSONLogging(getLevel(zap.InfoLevel))
	case debug:
		logger, err = initDevLogging(getLevel(zap.DebugLevel))
	default:
		logger, err = initSimpleLogging(getLevel(zap.DebugLevel))
	}
	if err != nil {
		return nil, errors.Wrap(err, "error initializing logging")
	}
	return logger.Sugar(), nil
}

func initDevLogging(lvl zapcore.Level) (*zap.Logger, error) {
	err := zap.RegisterEncoder(keyCustom, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewEncoder(cfg), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "error registering dev encoder")
	}
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zapcore.EncoderConfig{}
	config.Level = zap.NewAtomicLevelAt(lvl)
	config.Encoding = keyCustom
	return config.Build(zap.AddStacktrace(zap.PanicLevel), zap.AddCaller())
}

func initJSONLogging(lvl zapcore.Level) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.FunctionKey = "func"
	config.Level = zap.NewAtomicLevelAt(lvl)
	return config.Build(zap.AddStacktrace(zap.PanicLevel), zap.AddCaller())
}

func initSimpleLogging(lvl zapcore.Level) (*zap.Logger, error) {
	err := zap.RegisterEncoder(keyCustom, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return SimpleEncoder(cfg), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "error registering simple encoder")
	}
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zapcore.EncoderConfig{}
	config.Level = zap.NewAtomicLevelAt(lvl)
	config.Encoding = keyCustom
	return config.Build(zap.AddStacktrace(zap.PanicLevel), zap.AddCaller())
}

func getLevel(dflt zapcore.Level) zapcore.Level {
	l := os.Getenv("LOGGING_LEVEL")
	switch {
	case strings.EqualFold(l, "debug"):
		return zap.DebugLevel
	case strings.EqualFold(l, "info"):
		return zap.InfoLevel
	case strings.EqualFold(l, "warn"):
		return zap.WarnLevel
	case strings.EqualFold(l, "error"):
		return zap.ErrorLevel
	}
	return dflt
}
