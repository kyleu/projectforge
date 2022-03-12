package util

import (
	"strings"

	"github.com/pkg/errors"
)

func MarkdownTable(header []string, rows [][]string) (string, error) {
	maxes := make([]int, 0, len(header))
	for _, h := range header {
		maxes = append(maxes, len(h))
	}

	for vi, v := range rows {
		if len(v) != len(header) {
			return "", errors.Errorf("row [%d] contains [%d] fields, but [%d] header fields were provided", vi, len(v), len(header))
		}
		for cellIdx, cellVal := range v {
			if cellIdx <= len(maxes)-1 && len(cellVal) > maxes[cellIdx] {
				maxes[cellIdx] = len(cellVal)
			}
		}
	}

	ret := make([]string, 0, len(rows))
	add := func(x []string) {
		line := "| "
		for vi, v := range x {
			mx := 0
			if vi < len(maxes) {
				mx = maxes[vi]
			}
			line += StringPad(v, mx)
			if vi < len(x)-1 {
				line += " | "
			}
		}
		line += " |"
		ret = append(ret, line)
	}

	add(header)

	divider := "|-"
	for mi, m := range maxes {
		for i := 0; i < m; i++ {
			divider += "-"
		}
		if mi < len(maxes)-1 {
			divider += "-|-"
		}
	}
	divider += "-|"
	ret = append(ret, divider)

	for _, row := range rows {
		add(row)
	}
	return strings.Join(ret, "\n"), nil
}
