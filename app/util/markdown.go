package util

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func MarkdownTable(header []string, rows [][]string, linebreak string) (string, error) {
	maxes := lo.Map(header, func(h string, _ int) int {
		return len(h)
	})
	for vi, v := range rows {
		if len(v) != len(header) {
			return "", errors.Errorf("row [%d] contains [%d] fields, but [%d] header fields were provided", vi, len(v), len(header))
		}
		lo.ForEach(v, func(cellVal string, cellIdx int) {
			if cellIdx <= len(maxes)-1 && len(cellVal) > maxes[cellIdx] {
				maxes[cellIdx] = len(cellVal)
			}
		})
	}

	ret := make([]string, 0, len(rows))
	add := func(x []string) {
		line := "| "
		lo.ForEach(x, func(v string, vi int) {
			mx := 0
			if vi < len(maxes) {
				mx = maxes[vi]
			}
			line += StringPad(v, mx)
			if vi < len(x)-1 {
				line += " | "
			}
		})
		line += " |"
		ret = append(ret, line)
	}

	add(header)

	divider := "|-"
	lo.ForEach(maxes, func(m int, mi int) {
		lo.Times(m, func(_ int) struct{} {
			divider += "-"
			return struct{}{}
		})
		if mi < len(maxes)-1 {
			divider += "-|-"
		}
	})
	divider += "-|"
	ret = append(ret, divider)

	lo.ForEach(rows, func(row []string, _ int) {
		add(row)
	})
	return strings.Join(ret, linebreak), nil
}
