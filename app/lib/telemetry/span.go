// Content managed by Project Forge, see [projectforge.md] for details.
package telemetry

import (
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Span struct {
	OT trace.Span
}

func (s *Span) TraceID() string {
	return s.OT.SpanContext().TraceID().String()
}

func (s *Span) SpanID() string {
	return s.OT.SpanContext().SpanID().String()
}

func (s *Span) SetName(name string) {
	s.OT.SetName(name)
}

func (s *Span) SetStatus(status string, description string) {
	switch strings.ToLower(status) {
	case "ok":
		s.OT.SetStatus(codes.Ok, description)
	case "error":
		s.OT.SetStatus(codes.Error, description)
	default:
		s.OT.SetStatus(codes.Unset, description)
	}
}

func (s *Span) Attribute(k string, v interface{}) {
	s.Attributes(&Attribute{Key: k, Value: v})
}

func (s *Span) Attributes(attrs ...*Attribute) {
	ot := make([]attribute.KeyValue, 0, len(attrs))
	for _, attr := range attrs {
		ot = append(ot, attr.ToOT())
	}
	s.OT.SetAttributes(ot...)
}

func (s *Span) Event(name string, attrs ...*Attribute) {
	s.OT.AddEvent(name)
}

func (s *Span) OnError(err error) {
	s.OT.RecordError(err)
}

// Complete must be called, usually through a `defer` block.
func (s *Span) Complete() {
	s.OT.End()
}

func (s *Span) SetHTTPStatus(code int) {
	s.Attribute("http.status_code", code)
	x, desc := semconv.SpanStatusFromHTTPStatusCode(code)
	s.SetStatus(x.String(), desc)
}
