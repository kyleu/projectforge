package log

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type simpleEncoder struct {
	zapcore.Encoder
	pool buffer.Pool
}

func createSimpleEncoder(cfg zapcore.EncoderConfig) *simpleEncoder {
	return &simpleEncoder{Encoder: zapcore.NewJSONEncoder(cfg), pool: buffer.NewPool()}
}

func (e *simpleEncoder) Clone() zapcore.Encoder {
	return &simpleEncoder{Encoder: e.Encoder.Clone(), pool: e.pool}
}

func (e *simpleEncoder) EncodeEntry(entry zapcore.Entry, _ []zapcore.Field) (*buffer.Buffer, error) {
	ret := e.pool.Get()
	m := levelToColor[entry.Level.String()].Add(entry.Message)
	ret.AppendString(m)
	ret.AppendByte('\n')
	return ret, nil
}
