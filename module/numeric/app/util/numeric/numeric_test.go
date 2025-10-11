package numeric_test

import (
	"math"
	"testing"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util/numeric"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		test string
		res  numeric.Numeric
		err  error
	}{
		{"Invalid", "InvalidNum", numeric.Zero, errors.New("invalid number: InvalidNum")},
		{"NaN", "NaN", numeric.NaN, nil},
		{"Zero", "0", numeric.Zero, nil},
		{"One", "1", numeric.One, nil},
		{"Ten", "10", numeric.From(1, 1), nil},
		{"TenString", "10", numeric.Ten, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := numeric.FromString(tt.test)
			if err != nil {
				if tt.err != nil && tt.err.Error() != err.Error() {
					t.Error(err)
				}
				return
			}
			if !result.Equals(tt.res) && !math.IsNaN(tt.res.Mantissa()) {
				t.Errorf("numeric.FromString(%s) = %s; want %s", tt.test, result.Debug(), tt.res.Debug())
			}
		})
	}
}

type OpTest struct {
	name string
	l    numeric.Numeric
	r    numeric.Numeric
	res  numeric.Numeric
	err  error
}

func TestAdd(t *testing.T) {
	t.Parallel()
	tests := []*OpTest{
		{"Zero", numeric.Zero, numeric.Zero, numeric.Zero, nil},
		{"One", numeric.Zero, numeric.One, numeric.One, nil},
		{"Two", numeric.One, numeric.One, numeric.FromFloat(2), nil},
		{"Hundreds", numeric.FromFloat(123), numeric.FromFloat(321), numeric.FromFloat(444), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.l.Add(tt.r)
			if !result.Equals(tt.res) && !math.IsNaN(tt.res.Mantissa()) {
				t.Errorf("numeric.Add(%s, %s) = %s; want %s", tt.l, tt.r, result.Debug(), tt.res.Debug())
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	t.Parallel()
	tests := []*OpTest{
		{"Zero", numeric.Zero, numeric.Zero, numeric.Zero, nil},
		{"One", numeric.Zero, numeric.One, numeric.One.Negate(), nil},
		{"Two", numeric.One, numeric.One, numeric.Zero, nil},
		{"Hundreds", numeric.FromFloat(123), numeric.FromFloat(321), numeric.FromFloat(-198), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.l.Subtract(tt.r)
			if !result.Equals(tt.res) && !math.IsNaN(tt.res.Mantissa()) {
				t.Errorf("numeric.Subtract(%s, %s) = %s; want %s", tt.l, tt.r, result.Debug(), tt.res.Debug())
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	t.Parallel()
	tests := []*OpTest{
		{"Zero", numeric.Zero, numeric.Zero, numeric.Zero, nil},
		{"One", numeric.Zero, numeric.One, numeric.Zero, nil},
		{"Two", numeric.One, numeric.One, numeric.One, nil},
		{"Hundreds", numeric.FromFloat(123), numeric.FromFloat(321), numeric.FromFloat(39483), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.l.Multiply(tt.r).Round()
			if !result.Equals(tt.res) && !math.IsNaN(tt.res.Mantissa()) {
				t.Errorf("numeric.Multiply(%s, %s) = %s; want %s", tt.l, tt.r, result.Debug(), tt.res.Debug())
			}
		})
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		num  numeric.Numeric
		str  string
		json string
	}{
		{"Zero", numeric.Zero, `0`, `0`},
		{"One", numeric.One, `1`, `1`},
		{"OneFull", numeric.One, `1`, `{"m": 1,"e": 0}`},
		{"Ten", numeric.Ten, `10`, `10`},
		{"Negative", numeric.One.Negate(), `-1`, `-1`},
		{"Float", numeric.FromFloat(123.45), `123.45`, `123.45`},
		{"Big", numeric.From(1.01, 1000), `{"m": 1.01, "e": 1000}`, `{"m": 1.01, "e": 1000}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.num.MarshalJSON()
			if err != nil {
				t.Errorf("numeric.MarshalJSON(%s) error = %v", tt.num, err)
				return
			}
			if string(result) != tt.str {
				t.Errorf("numeric.MarshalJSON(%s) = %s; want %s", tt.num, string(result), tt.str)
			}

			var unmarshaled numeric.Numeric
			err = unmarshaled.UnmarshalJSON([]byte(tt.json))
			if err != nil {
				t.Errorf("numeric.UnmarshalJSON(%s) error = %v", tt.json, err)
				return
			}
			if !unmarshaled.Equals(tt.num) {
				t.Errorf("numeric.UnmarshalJSON(%s) = %s; want %s", tt.str, unmarshaled.Debug(), tt.num.Debug())
			}
		})
	}
}
