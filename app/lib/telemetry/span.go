// Content managed by Project Forge, see [projectforge.md] for details.
package telemetry

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"projectforge.dev/projectforge/app/util"
)

const disabledMsg = "[telemetry.disabled]"

type Span struct {
	OT        trace.Span
	statusSet bool
}

func (s *Span) TraceID() string {
	if s == nil || !enabled {
		return disabledMsg
	}
	return s.OT.SpanContext().TraceID().String()
}

func (s *Span) SpanID() string {
	if s == nil || !enabled {
		return disabledMsg
	}
	return s.OT.SpanContext().SpanID().String()
}

func (s *Span) SetName(name string) {
	if s == nil || !enabled {
		return
	}
	s.OT.SetName(name)
}

func (s *Span) SetStatus(status string, description string) {
	if s == nil || !enabled {
		return
	}
	s.statusSet = true
	switch strings.ToLower(status) {
	case "ok":
		s.OT.SetStatus(codes.Ok, description)
	case util.KeyError:
		s.OT.SetStatus(codes.Error, description)
	default:
		s.OT.SetStatus(codes.Ok, status+": "+description)
	}
}

func (s *Span) Attribute(k string, v any) {
	if s == nil || !enabled {
		return
	}
	s.Attributes(&Attribute{Key: k, Value: v})
}

func (s *Span) Attributes(attrs ...*Attribute) {
	if s == nil || !enabled {
		return
	}
	ot := lo.Map(attrs, func(attr *Attribute, _ int) attribute.KeyValue {
		return attr.ToOT()
	})
	s.OT.SetAttributes(ot...)
}

func (s *Span) Event(name string, attrs ...*Attribute) {
	if s == nil || !enabled {
		return
	}
	s.OT.AddEvent(name)
	lo.ForEach(attrs, func(attr *Attribute, _ int) {
		s.OT.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key(attr.Key),
			Value: attribute.StringValue(fmt.Sprint(attr.Value)),
		})
	})
}

func (s *Span) OnError(err error) {
	if s == nil || !enabled {
		return
	}
	s.OT.RecordError(err)
}

// Complete must be called, usually through a `defer` block.
func (s *Span) Complete() {
	if s == nil || !enabled {
		return
	}
	if !s.statusSet {
		s.SetStatus("ok", "complete")
	}
	s.OT.End()
}

func (s *Span) SetHTTPStatus(code int) {
	if s == nil || !enabled {
		return
	}
	s.Attribute("http.status_code", code)
	x, desc := semconv.SpanStatusFromHTTPStatusCode(code)
	s.SetStatus(x.String(), desc)
}

func (s *Span) String() string {
	if s == nil || !enabled {
		return disabledMsg
	}
	return s.SpanID() + "::" + s.TraceID()
}
