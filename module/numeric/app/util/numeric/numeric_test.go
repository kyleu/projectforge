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
