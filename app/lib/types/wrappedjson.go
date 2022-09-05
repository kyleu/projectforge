// Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"encoding/json"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

type wrappedUnmarshal struct {
	K string          `json:"k"`
	T json.RawMessage `json:"t,omitempty"`
}

func (x *Wrapped) MarshalJSON() ([]byte, error) {
	b := util.ToJSONBytes(x.T, false)
	// needs better detection
	if len(b) == 2 {
		return util.ToJSONBytes(x.K, false), nil
	}
	return util.ToJSONBytes(wrappedUnmarshal{K: x.K, T: b}, false), nil
}

//nolint:funlen, gocyclo, cyclop
func (x *Wrapped) UnmarshalJSON(data []byte) error {
	var wu wrappedUnmarshal
	err := util.FromJSON(data, &wu)
	if err != nil {
		str := ""
		newErr := util.FromJSON(data, &str)
		if newErr != nil {
			return err
		}
		wu = wrappedUnmarshal{K: str, T: []byte("{}")}
	}
	var t Type
	switch wu.K {
	case KeyAny:
		tgt := &Any{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyBit:
		tgt := &Bit{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyBool:
		tgt := &Bool{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyByte:
		tgt := &Byte{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyChar:
		tgt := &Char{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyDate:
		tgt := &Date{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyEnum:
		tgt := &Enum{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyEnumValue:
		tgt := &EnumValue{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyError:
		tgt := &Error{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyFloat:
		tgt := &Float{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyInt:
		tgt := &Int{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyJSON:
		tgt := &JSON{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyList:
		tgt := &List{}
		err = util.FromJSON(wu.T, &tgt)
		if tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyMap:
		tgt := &Map{}
		err = util.FromJSON(wu.T, &tgt)
		if tgt.K == nil {
			tgt.K = NewString()
		}
		if tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyMethod:
		tgt := &Method{}
		err = util.FromJSON(wu.T, &tgt)
		if tgt.Ret == nil {
			tgt.Ret = NewAny()
		}
		t = tgt
	case KeyNil:
		tgt := &Nil{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyOption:
		tgt := &Option{}
		err = util.FromJSON(wu.T, &tgt)
		if tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyRange:
		tgt := &Range{}
		err = util.FromJSON(wu.T, &tgt)
		if tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyReference:
		tgt := &Reference{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeySet:
		tgt := &Set{}
		err = util.FromJSON(wu.T, &tgt)
		if tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyString:
		tgt := &String{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyTime:
		tgt := &Time{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyTimestamp:
		tgt := &Timestamp{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyTimestampZoned:
		tgt := &TimestampZoned{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyUnknown:
		tgt := &Unknown{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyUUID:
		tgt := &UUID{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyValueMap:
		tgt := &ValueMap{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	case KeyXML:
		tgt := &XML{}
		err = util.FromJSON(wu.T, &tgt)
		t = tgt
	default:
		t = &Unknown{X: "unmarshal:" + wu.K}
	}
	if err != nil {
		return errors.Wrapf(err, "unable to unmarshal wrapped field of type [%s]", wu.K)
	}
	if t == nil {
		return errors.New("nil type returned from unmarshal")
	}
	x.K = wu.K
	x.T = t
	return nil
}
