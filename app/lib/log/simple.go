package log

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"

	"projectforge.dev/projectforge/app/util"
)

type simpleEncoder struct {
	zapcore.Encoder
	pool      buffer.Pool
	listeners []ListenerFunc
}

func createSimpleEncoder(cfg zapcore.EncoderConfig, fns ...ListenerFunc) *simpleEncoder {
	return &simpleEncoder{Encoder: zapcore.NewJSONEncoder(cfg), pool: buffer.NewPool(), listeners: fns}
}

func (e *simpleEncoder) Clone() zapcore.Encoder {
	return &simpleEncoder{Encoder: e.Encoder.Clone(), pool: e.pool}
}

func (e *simpleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	go func() {
		e.sendToListeners(entry, fields)
	}()
	ret := e.pool.Get()
	m := levelToColor[entry.Level.String()].Add(entry.Message)
	ret.AppendString(m)
	ret.AppendByte('\n')
	return ret, nil
}

func (e *simpleEncoder) sendToListeners(entry zapcore.Entry, fields []zapcore.Field) {
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
