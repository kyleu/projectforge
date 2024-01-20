package gql

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func All(models model.Models, enums enum.Enums, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"app", "gql"}, "generated.graphql")
	headerBlock := golang.NewBlock("header", "graphql")
	headerBlock.W("### This file contains generated definitions of generated models and enums")
	headerBlock.W("### To use these definitions, include them inside your `schema.graphql` file")
	g.AddBlocks(headerBlock)
	if len(enums) > 0 {
		enumBlocks(g, enums)
	}
	if len(enums) > 0 {
		modelBlocks(g, models, enums)
	}
	return g.Render(addHeader, linebreak)
}

func enumBlocks(g *golang.Template, enums enum.Enums) {
	for _, e := range enums {
		g.AddBlocks(enumBlock(e))
	}
}

func enumBlock(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock("enum:"+e.Name, "graphql")
	ret.WB()
	ret.W("# Enum Type [%s]", e.Proper())
	return ret
}

func modelBlocks(g *golang.Template, models model.Models, enums enum.Enums) {
	for _, m := range models {
		g.AddBlocks(modelBlock(m, enums))
	}
}

func modelBlock(m *model.Model, _ enum.Enums) *golang.Block {
	ret := golang.NewBlock("model:"+m.Name, "graphql")
	ret.WB()
	ret.W("# Model Class [%s]", m.Title())
	return ret
}
