package golang

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/util"
)

type Block struct {
	Key   string   `json:"key"`
	Type  string   `json:"type"`
	Lines []string `json:"lines"`
}

func NewBlock(k string, t string) *Block {
	return &Block{Key: k, Type: t}
}

func (b *Block) W(l string, args ...any) {
	if len(args) == 0 {
		b.Lines = append(b.Lines, l)
	} else {
		b.Lines = append(b.Lines, fmt.Sprintf(l, args...))
	}
}

func (b *Block) WB() {
	b.Lines = append(b.Lines, "")
}

func (b *Block) WE(indent int) {
	ind := util.StringRepeat("\t", indent)
	b.Lines = append(b.Lines, "%sif err != nil {", ind)
	b.Lines = append(b.Lines, "%s\t return err", ind)
	b.Lines = append(b.Lines, "%s}", ind)
}

func (b *Block) Render() string {
	return strings.Join(b.Lines, "\n")
}

type Blocks []*Block
