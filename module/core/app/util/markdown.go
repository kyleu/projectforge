package util

import (
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

	ret := NewStringSlice(make([]string, 0, len(rows)))
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
		ret.Push(line)
	}

	add(header)

	divider := "|-"
	lo.ForEach(maxes, func(m int, mi int) {
		lo.Times(m, func(_ int) struct{} {
			divider += "-"
			return EmptyStruct
		})
		if mi < len(maxes)-1 {
			divider += "-|-"
		}
	})
	divider += "-|"
	ret.Push(divider)

	lo.ForEach(rows, func(row []string, _ int) {
		add(row)
	})
	return ret.Join(linebreak), nil
}

func MarkdownTableParse(md string) ([]string, [][]string) {
	var header []string
	var rows [][]string
	for idx, line := range StringSplitLines(md) {
		split := StringSplitAndTrim(line, "|")
		if idx == 0 {
			header = split
		} else if idx != 1 {
			rows = append(rows, split)
		}
	}
	return header, rows
}
