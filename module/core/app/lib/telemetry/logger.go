package telemetry

import (
	"strings"

	"go.uber.org/zap"

	"{{{ .Package }}}/app/util"
)

func LoggerFor(logger util.Logger, span *Span) util.Logger {
	if span == nil {
		return logger
	}
	return logger.With(zap.String("trace", span.TraceID()), zap.String("span", span.SpanID()))
}

type ErrHandler struct {
	logger     util.Logger
	hasPrinted bool
}

func (e *ErrHandler) Handle(err error) {
	if err == nil {
		return
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "Post \"") {
		if e.hasPrinted {
			return
		}
		if idx := strings.Index(msg, "\":"); idx > -1 {
			msg = strings.TrimSpace(msg[idx+2:])
		}
		e.logger.Warn("telemetry seems to be unavailable: [" + msg + "] (this message will appear only once)")
		e.hasPrinted = true
		return
	}
	e.logger.Warnf("telemetry error: %+v", err)
}
