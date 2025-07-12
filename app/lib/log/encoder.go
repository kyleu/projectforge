package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"

	"projectforge.dev/projectforge/app/util"
)

const (
	timeFormat = "15:04:05.000000"
	logIndent  = "  "
)

type ListenerFunc func(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap)

type customEncoder struct {
	zapcore.Encoder
	colored   bool
	pool      buffer.Pool
	listeners []ListenerFunc
}

func newEncoder(cfg zapcore.EncoderConfig, colored bool, fns ...ListenerFunc) *customEncoder {
	return &customEncoder{Encoder: zapcore.NewJSONEncoder(cfg), colored: colored, pool: buffer.NewPool(), listeners: fns}
}

func (e *customEncoder) Clone() zapcore.Encoder {
	return &customEncoder{Encoder: e.Encoder.Clone(), colored: e.colored, pool: e.pool}
}

func (e *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	go func() {
		recentMU.Lock()
		defer recentMU.Unlock()
		RecentLogs = append(RecentLogs, &entry)
		if len(RecentLogs) > 50 {
			RecentLogs = RecentLogs[1:]
		}
	}()
	go func() {
		e.sendToListeners(entry, fields)
	}()
	b, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, errors.Wrap(err, "logging error")
	}
	out := b.Bytes()
	b.Free()

	data, err := util.FromJSONMap(out)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse logging JSON")
	}

	ret := e.pool.Get()
	ret.AppendByte('\n')
	addLine := func(l string) {
		ret.AppendString(l)
		ret.AppendByte('\n')
	}

	lvl := fmt.Sprintf("%-5v", entry.Level.CapitalString())
	if e.colored {
		lvl = levelToColor[entry.Level.String()].Add(lvl)
	}
	tm := entry.Time.Format(timeFormat)

	msg := entry.Message
	var msgLines []string
	if strings.Contains(msg, "\n") {
		msgLines = util.StringSplitLines(msg)
		msg = msgLines[0]
		msgLines = msgLines[1:]
	}

	if e.colored {
		addLine(fmt.Sprintf("[%s] %s %s", lvl, tm, Cyan.Add(msg)))
	} else {
		addLine(fmt.Sprintf("[%s] %s %s", lvl, tm, msg))
	}

	lo.ForEach(msgLines, func(ml string, _ int) {
		if e.colored {
			if strings.Contains(ml, util.AppKey) {
				ml = Green.Add(ml)
			}
			addLine(logIndent + Cyan.Add(ml))
		} else {
			addLine(logIndent + ml)
		}
	})
	if len(data) > 0 {
		addLine(logIndent + util.ToJSONCompact(data))
	}
	caller := entry.Caller.String()
	if idx := strings.Index(caller, "/app"); idx > 0 {
		caller = caller[idx+1:]
	}
	if entry.Caller.Function != "" {
		caller += " (" + entry.Caller.Function + ")"
	}
	addLine(logIndent + caller)

	if entry.Stack != "" {
		st := util.StringSplitLines(entry.Stack)
		lo.ForEach(st, func(stl string, _ int) {
			if strings.Contains(stl, util.AppKey) {
				stl = Green.Add(stl)
			}
			addLine(logIndent + stl)
		})
	}
	return ret, nil
}

func (e *customEncoder) sendToListeners(entry zapcore.Entry, fields []zapcore.Field) {
	listenerMU.Lock()
	defer listenerMU.Unlock()
	fieldMap := make(util.ValueMap, len(fields))
	for _, x := range fields {
		fieldMap[x.Key] = x.Interface
	}
	caller := util.ValueMap{"file": entry.Caller.File, "line": entry.Caller.Line, "function": entry.Caller.Function}
	for _, listener := range e.listeners {
		l := listener
		go func() {
			l(entry.Level.String(), entry.Time, entry.LoggerName, entry.Message, caller, entry.Stack, fieldMap)
		}()
	}
}
