package svc

import "fmt"

type Model interface {
	fmt.Stringer
	TitleString() string
	Strings() []string
	ToCSV() ([]string, [][]string)
	ToData() []any
}

type ModelSeq []Model
