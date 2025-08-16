package numeric

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

type NumericJSON struct {
	M float64 `json:"m"`
	E int64   `json:"e,omitempty"`
}

func (n *NumericJSON) ToNumeric() Numeric {
	return From(n.M, n.E)
}

func (n Numeric) MarshalJSON() ([]byte, error) {
	if n.LessThan(MaxString) {
		return []byte(n.String()), nil
	}
	return []byte(`{"m": ` + fmt.Sprint(n.mantissa) + `, "e": ` + fmt.Sprint(n.exponent) + `}`), nil
}

func (n *Numeric) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	l := len(data)
	if l == 0 || (l == 4 && slices.Equal(data, []byte("null"))) {
		return nil
	}
	switch data[0] {
	case '{':
		s, err := util.FromJSONObj[NumericJSON](data)
		if err != nil {
			return errors.Wrapf(err, "invalid object value [%s]", string(data))
		}
		*n = s.ToNumeric()
	case '"', '\'':
		s, err := util.FromJSONString(data)
		if err != nil {
			return errors.Wrapf(err, "invalid string value [%s]", string(data))
		}
		ret, err := FromString(s)
		if err != nil {
			return errors.Wrapf(err, "invalid numeric string [%s]", string(data))
		}
		*n = ret
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s, err := util.FromJSONObj[float64](data)
		if err != nil {
			return errors.Wrapf(err, "invalid float value [%s]", string(data))
		}
		ret := FromFloat(s)
		*n = ret
	default:
		return errors.Errorf("invalid numeric content [%s]", string(data))
	}
	return nil
}
