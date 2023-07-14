package telemetry

import (
	"strings"

	"{{{ .Package }}}/app/util"
)

func LoggerFor(logger util.Logger, span *Span) util.Logger {
	if logger == nil {
		return nil
	}
	if span == nil {
		return logger
	}
	return logger.With("trace", span.TraceID(), "span", span.SpanID())
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
	if strings.HasPrefix(msg, "Post \"") || strings.HasPrefix(msg, "traces export") {
		if e.hasPrinted {
			return
		}
		if idx := strings.Index(msg, "\":"); idx > -1 {
			msg = strings.TrimSpace(msg[idx+2:])
		}
		if strings.Contains(msg, "connectex") {
			msg = "connection failure"
		}
		e.logger.Warn("telemetry seems to be unavailable: [" + msg + "] (this message will appear only once)")
		e.hasPrinted = true
		return
	}
	e.logger.Warnf("telemetry error: %+v", err)
}
