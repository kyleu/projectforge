// Package util - Content managed by Project Forge, see [projectforge.md] for details.
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
	if x, ok := data.(csvInterface); ok {
		cols, rows := x.ToCSV()
		return cols, rows, nil
	}
	if x, ok := data.(string); ok {
		csvdata, err := csv.NewReader(strings.NewReader(x)).ReadAll()
		if err != nil || len(csvdata) == 0 {
			return []string{x}, [][]string{}, nil //nolint:nilerr
		}
		return csvdata[0], csvdata[1:], nil
	}
	if x, ok := data.(error); ok {
		return []string{"Error", "Message"}, [][]string{{fmt.Sprintf("%T", x), x.Error()}}, nil
	}
	if x, ok := data.(*ErrorDetail); ok {
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
