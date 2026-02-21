package mvc

import (
	"context"
	"fmt"
	"time"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/util"
)

type PageState struct {
	Action    string          `json:"action,omitzero"`
	Title     string          `json:"title,omitzero"`
	Cursor    int             `json:"cursor,omitzero"`
	Status    string          `json:"status,omitzero"`
	Error     string          `json:"error,omitzero"`
	Data      util.ValueMap   `json:"data,omitzero"`
	UpdatedAt time.Time       `json:"updatedAt,omitzero"`
	Logger    util.Logger     `json:"-"`
	Context   context.Context `json:"-"` //nolint:containedctx // owned per-page and closed via Close
	Span      *telemetry.Span `json:"-"`
}

func NewPageState(parentCtx context.Context, act string, title string, data util.ValueMap, logger util.Logger) *PageState {
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	ctx, span, scopedLogger := telemetry.StartSpan(parentCtx, "tui:"+act, logger)
	span.Attribute("action", act)
	return &PageState{Action: act, Title: title, Data: data, UpdatedAt: util.TimeCurrent(), Logger: scopedLogger, Context: ctx, Span: span}
}

func (ps *PageState) EnsureData() util.ValueMap {
	if ps.Data == nil {
		ps.Data = util.ValueMap{}
	}
	return ps.Data
}

func (ps *PageState) SetStatus(msg string, args ...any) {
	if len(args) > 0 {
		ps.Status = fmt.Sprintf(msg, args...)
	} else {
		ps.Status = msg
	}
	ps.Error = ""
	ps.UpdatedAt = util.TimeCurrent()
}

func (ps *PageState) SetStatusText(msg string) {
	ps.Status = msg
	ps.Error = ""
	ps.UpdatedAt = util.TimeCurrent()
}

func (ps *PageState) SetError(err error) {
	if err == nil {
		ps.Error = ""
		return
	}
	ps.Error = err.Error()
	ps.UpdatedAt = util.TimeCurrent()
}

func (ps *PageState) Close() {
	if ps != nil && ps.Span != nil {
		ps.Span.Complete()
	}
}
