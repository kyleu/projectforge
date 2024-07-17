package golang

import (
	"fmt"
	"go/parser"
	"go/token"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/uudashr/gocognit"

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

func (b *Block) WA(a ...string) {
	lo.ForEach(a, func(x string, _ int) {
		b.W(x)
	})
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

func (b *Block) Render(linebreak string) (string, error) {
	d, err := b.NoLineDecl(linebreak)
	if err != nil {
		return "", err
	}

	if d == "" {
		if len(b.Lints) == 0 {
			return strings.Join(b.Lines, linebreak), nil
		}
		line := "//nolint:" + strings.Join(b.Lints, ", ")
		return strings.Join(append([]string{line}, b.Lines...), linebreak), nil
	}
	line := strings.Join(append([]string{d}, b.Lints...), ", ")
	return strings.Join(append([]string{line}, b.Lines...), linebreak), nil
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

func (b *Block) NoLineDecl(linebreak string) (string, error) {
	ret := &util.StringSlice{}
	if b.LineCount() > 102 && strings.Contains(b.Lines[0], "func") {
		ret.Push("funlen")
	}
	if b.LineMaxLength() > 160 {
		ret.Push("lll")
	}

	fset := token.NewFileSet()
	if f, _ := parser.ParseFile(fset, "temp.go", strings.Join(append([]string{"package x"}, b.Lines...), linebreak), 0); f != nil {
		for _, complexity := range gocognit.ComplexityStats(f, fset, nil) {
			if complexity.Complexity > 30 && !slices.Contains(ret.Slice, "gocognit") {
				ret.Push("gocognit")
			}
		}
	}
	if len(ret.Slice) == 0 || b.SkipDecl || b.ContainsText("{%") {
		return "", nil
	}
	x := fmt.Sprintf("//nolint:%s", ret.Join(","))
	return x, nil
}

func (b *Block) ContainsText(s string) bool {
	return lo.ContainsBy(b.Lines, func(l string) bool {
		return strings.Contains(l, s)
	})
}

type Blocks []*Block
