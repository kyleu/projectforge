package types

import (
	"encoding/json"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

type wrappedUnmarshal struct {
	K string          `json:"k"`
	T json.RawMessage `json:"t,omitempty"`
}

func (x *Wrapped) MarshalJSON() ([]byte, error) {
	b := util.ToJSONBytes(x.T, false)
	s := string(b)
	if s == objStr || (x.T.Key() == KeyMap && s == `{"k":"string","v":"any"}`) || ((x.T.Key() == KeyList || x.T.Key() == KeyOrderedMap) && s == `{"v":"any"}`) {
		return util.ToJSONBytes(x.K, false), nil
	}
	return util.ToJSONBytes(wrappedUnmarshal{K: x.K, T: b}, false), nil
}

//nolint:funlen, gocyclo, cyclop
func (x *Wrapped) UnmarshalJSON(data []byte) error {
	var wu wrappedUnmarshal
	err := util.FromJSON(data, &wu)
	if err != nil {
		var str string
		newErr := util.FromJSON(data, &str)
		if newErr != nil {
			return err
		}
		wu = wrappedUnmarshal{K: str, T: []byte(objStr)}
	}
	if wu.K == "boolean" {
		wu.K = KeyBool
	}
	var t Type
	switch wu.K {
	case KeyAny:
		t, err = util.FromJSONObj[*Any](wu.T)
	case KeyBit:
		t, err = util.FromJSONObj[*Bit](wu.T)
	case KeyBool:
		t, err = util.FromJSONObj[*Bool](wu.T)
	case KeyByte:
		t, err = util.FromJSONObj[*Byte](wu.T)
	case KeyChar:
		t, err = util.FromJSONObj[*Char](wu.T)
	case KeyDate:
		t, err = util.FromJSONObj[*Date](wu.T)
	case KeyEnum:
		t, err = util.FromJSONObj[*Enum](wu.T)
	case KeyEnumValue:
		t, err = util.FromJSONObj[*EnumValue](wu.T)
	case KeyError:
		t, err = util.FromJSONObj[*Error](wu.T)
	case KeyFloat:
		t, err = util.FromJSONObj[*Float](wu.T)
	case KeyInt:
		t, err = util.FromJSONObj[*Int](wu.T)
	case KeyJSON:
		t, err = util.FromJSONObj[*JSON](wu.T)
	case KeyList:
		var tgt *List
		tgt, err = util.FromJSONObj[*List](wu.T)
		if err == nil && tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyMap:
		var tgt *Map
		tgt, err = util.FromJSONObj[*Map](wu.T)
		if err == nil && tgt.K == nil {
			tgt.K = NewString()
		}
		if err == nil && tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyMethod:
		var tgt *Method
		tgt, err = util.FromJSONObj[*Method](wu.T)
		if tgt.Ret == nil {
			tgt.Ret = NewAny()
		}
		t = tgt
	case KeyNil:
		t, err = util.FromJSONObj[*Nil](wu.T)
	case KeyNumeric:
		t, err = util.FromJSONObj[*Numeric](wu.T)
	case KeyOption:
		var tgt *Option
		tgt, err = util.FromJSONObj[*Option](wu.T)
		if tgt != nil && tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyRange:
		var tgt *Range
		tgt, err = util.FromJSONObj[*Range](wu.T)
		if tgt != nil && tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyReference:
		t, err = util.FromJSONObj[*Reference](wu.T)
	case KeySet:
		var tgt *Set
		tgt, err = util.FromJSONObj[*Set](wu.T)
		if tgt != nil && tgt.V == nil {
			tgt.V = NewAny()
		}
		t = tgt
	case KeyString:
		t, err = util.FromJSONObj[*String](wu.T)
	case KeyTime:
		t, err = util.FromJSONObj[*Time](wu.T)
	case KeyTimestamp:
		t, err = util.FromJSONObj[*Timestamp](wu.T)
	case KeyTimestampZoned:
		t, err = util.FromJSONObj[*TimestampZoned](wu.T)
	case KeyUnknown:
		t, err = util.FromJSONObj[*Unknown](wu.T)
	case KeyUUID:
		t, err = util.FromJSONObj[*UUID](wu.T)
	case KeyValueMap:
		t, err = util.FromJSONObj[*ValueMap](wu.T)
	case KeyXML:
		t, err = util.FromJSONObj[*XML](wu.T)
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
