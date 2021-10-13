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
		logger, err = initJSONLogging()
	case debug:
		logger, err = initDevLogging()
	default:
		logger, err = initSimpleLogging()
	}
	if err != nil {
		return nil, errors.Wrap(err, "error initializing logging")
	}
	return logger.Sugar(), nil
}

func initDevLogging() (*zap.Logger, error) {
	err := zap.RegisterEncoder(keyCustom, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewEncoder(cfg), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "error registering dev encoder")
	}
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zapcore.EncoderConfig{}

	config.Encoding = keyCustom
	return config.Build(zap.AddStacktrace(zap.PanicLevel), zap.AddCaller())
}

func initJSONLogging() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.FunctionKey = "func"
	return config.Build(zap.AddStacktrace(zap.PanicLevel), zap.AddCaller())
}

func initSimpleLogging() (*zap.Logger, error) {
	err := zap.RegisterEncoder(keyCustom, func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return SimpleEncoder(cfg), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "error registering simple encoder")
	}
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zapcore.EncoderConfig{}

	config.Encoding = keyCustom
	return config.Build(zap.AddStacktrace(zap.PanicLevel), zap.AddCaller())
}
