package log

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"

	"{{{ .Package }}}/app/util"
)

const timeFormat = "15:04:05.000000"

type customEncoder struct {
	zapcore.Encoder
	pool buffer.Pool
}

// nolint
func NewEncoder(cfg zapcore.EncoderConfig) *customEncoder {
	return &customEncoder{Encoder: zapcore.NewJSONEncoder(cfg), pool: buffer.NewPool()}
}

func (e *customEncoder) Clone() zapcore.Encoder {
	return &customEncoder{Encoder: e.Encoder.Clone(), pool: e.pool}
}

// nolint
func (e *customEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	b, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, errors.Wrap(err, "logging error")
	}
	out := b.Bytes()
	b.Free()

	data := map[string]interface{}{}
	err = util.FromJSON(out, &data)
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
	lvl = levelToColor[entry.Level].Add(lvl)
	tm := entry.Time.Format(timeFormat)

	msg := entry.Message
	var msgLines []string
	if strings.Contains(msg, "\n") {
		msgLines = strings.Split(msg, "\n")
		msg = msgLines[0]
		msgLines = msgLines[1:]
	}

	addLine(fmt.Sprintf("[%s] %s %s", lvl, tm, Cyan.Add(msg)))
	for _, ml := range msgLines {
		if strings.Contains(ml, util.AppKey) {
			ml = Green.Add(ml)
		}
		addLine("  " + Cyan.Add(ml))
	}
	if len(data) > 0 {
		addLine("  " + util.ToJSONCompact(data))
	}
	caller := entry.Caller.String()
	if entry.Caller.Function != "" {
		caller += " (" + entry.Caller.Function + ")"
	}
	// idx := strings.Index(caller, "github.com/")
	// if idx > 0 {
	// 	caller = caller[idx:]
	// }
	addLine("  " + caller)

	if entry.Stack != "" {
		st := strings.Split(entry.Stack, "\n")
		for _, stl := range st {
			if strings.Contains(stl, util.AppKey) {
				stl = Green.Add(stl)
			}
			addLine("  " + stl)
		}
	}
	return ret, nil
}
