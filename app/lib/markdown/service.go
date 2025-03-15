package markdown

import (
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Service struct {
	Parser parser.Parser
}

func NewService() *Service {
	p := goldmark.DefaultParser()
	return &Service{Parser: p}
}

func (s *Service) Parse(data []byte) ast.Node {
	reader := text.NewReader(data)
	return s.Parser.Parse(reader)
}

func (s *Service) HTML(n ast.Node, data []byte, depth int) string {
	switch n.Type() {
	case ast.TypeDocument:
		return fmt.Sprintf("Document: %s", string(n.Text(data)))
	case ast.TypeInline:
		return fmt.Sprintf("Inline: %s", string(n.Text(data)))
	case ast.TypeBlock:
		return fmt.Sprintf("Block: %s", string(n.Text(data)))
	default:
		return fmt.Sprintf("unknown node type [%v]", n.Type())
	}
}
