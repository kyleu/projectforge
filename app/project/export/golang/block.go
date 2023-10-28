package golang

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Block struct {
	Key      string   `json:"key"`
	Type     string   `json:"type"`
	Lines    []string `json:"lines"`
	SkipDecl bool     `json:"skipDecl,omitempty"`
	Lints    []string `json:"lints,omitempty"`
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

func (b *Block) WE(indent int, prefix ...string) {
	ind := util.StringRepeat("\t", indent)
	p := strings.Join(prefix, ", ")
	if p != "" {
		p += ", "
	}
	b.Lines = append(b.Lines,
		ind+"if err != nil {",
		ind+"\treturn "+p+"err",
		ind+"}",
	)
}

func (b *Block) Render(linebreak string) string {
	d := b.NoLineDecl()
	if d == "" {
		if len(b.Lints) == 0 {
			return strings.Join(b.Lines, linebreak)
		}
		line := "//nolint:" + strings.Join(b.Lints, ", ")
		return strings.Join(append([]string{line}, b.Lines...), linebreak)
	}
	line := strings.Join(append([]string{d}, b.Lints...), ", ")
	return strings.Join(append([]string{line}, b.Lines...), linebreak)
}

func (b *Block) LineCount() int {
	return len(b.Lines)
}

func (b *Block) LineMaxLength() int {
	return util.StringArrayMaxLength(b.Lines)
}

func (b *Block) LineComplexity() int {
	return lo.SumBy(b.Lines, func(l string) int {
		tests := []string{"if ", "else ", "switch ", "select ", "for ", "break", "continue", "&&", "||"}
		return lo.Sum(lo.Map(tests, func(x string, _ int) int {
			return strings.Count(l, x)
		}))
	})
}

func (b *Block) NoLineDecl() string {
	var ret []string
	if b.LineCount() > 100 {
		ret = append(ret, "funlen")
	}
	if b.LineMaxLength() > 160 {
		ret = append(ret, "lll")
	}
	if b.LineComplexity() >= 28 {
		ret = append(ret, "gocognit")
	}
	if len(ret) == 0 || b.SkipDecl || b.ContainsText("{%") {
		return ""
	}
	x := fmt.Sprintf("//nolint:%s", strings.Join(ret, ","))
	return x
}

func (b *Block) ContainsText(s string) bool {
	return lo.ContainsBy(b.Lines, func(l string) bool {
		return strings.Contains(l, s)
	})
}

type Blocks []*Block
