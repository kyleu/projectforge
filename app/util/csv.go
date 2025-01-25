package util

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type csvInterface interface {
	ToCSV() ([]string, [][]string)
}

func ToCSV(data any) ([]string, [][]string, error) {
	if x, err := Cast[csvInterface](data); err == nil {
		cols, rows := x.ToCSV()
		return cols, rows, nil
	}
	if x, err := Cast[string](data); err == nil {
		csvdata, err := csv.NewReader(strings.NewReader(x)).ReadAll()
		if err != nil || len(csvdata) == 0 {
			return []string{x}, [][]string{}, nil //nolint:nilerr
		}
		return csvdata[0], csvdata[1:], nil
	}
	if x, err := Cast[error](data); err == nil {
		return []string{"Error", "Message"}, [][]string{{fmt.Sprintf("%T", x), x.Error()}}, nil
	}
	if x, err := Cast[*ErrorDetail](data); err == nil {
		return []string{"Error", "Message"}, [][]string{{x.Type, x.Message}}, nil
	}
	return nil, nil, errors.Errorf("unable to transform [%T] to CSV", data)
}

func ToCSVBytes(data any) ([]byte, error) {
	cols, rows, err := ToCSV(data)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(nil)
	c := csv.NewWriter(b)
	err = c.Write(cols)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		err = c.Write(row)
		if err != nil {
			return nil, err
		}
	}
	c.Flush()
	return b.Bytes(), nil
}
