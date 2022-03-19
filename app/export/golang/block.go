package golang

import (
	"fmt"
	"strings"
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

func (b *Block) Render() string {
	return strings.Join(b.Lines, "\n")
}

type Blocks []*Block
