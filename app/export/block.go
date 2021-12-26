package export

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

func (b *Block) Priority() int {
	switch b.Type {
	default:
		return 0
	}
}

func (b *Block) W(l string) {
	b.Lines = append(b.Lines, l)
}

func (b *Block) WF(l string, args ...interface{}) {
	b.W(fmt.Sprintf(l, args...))
}

func (b *Block) Render() string {
	return strings.Join(b.Lines, "\n")
}

type Blocks []*Block
