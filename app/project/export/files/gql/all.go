package gql

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func All(models model.Models, enums enum.Enums, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"app", "gql"}, "generated.graphql")
	if len(enums) > 0 {
		enumBlocks(g, enums)
	}
	if len(models) > 0 {
		modelBlocks(g, models, enums)
	}
	g.AddBlocks(extendQuery(), extendMutation())
	return g.Render(linebreak)
}

func enumBlocks(g *golang.Template, enums enum.Enums) {
	for _, e := range enums {
		g.AddBlocks(enumBlock(e))
	}
}

func enumBlock(e *enum.Enum) *golang.Block {
	ret := golang.NewBlock("enum:"+e.Name, "graphql")
	ret.WF("# Enum Type [%s]", e.Proper())
	return ret
}

func modelBlocks(g *golang.Template, models model.Models, enums enum.Enums) {
	for _, m := range models {
		g.AddBlocks(modelBlock(m, enums))
	}
}

func modelBlock(m *model.Model, _ enum.Enums) *golang.Block {
	ret := golang.NewBlock("model:"+m.Name, "graphql")
	ret.WF("# Model Class [%s]", m.Title())
	return ret
}

func extendQuery() *golang.Block {
	ret := golang.NewBlock("extendQuery", "graphql")
	ret.W("extend type Query {")
	ret.W("}")
	return ret
}

func extendMutation() *golang.Block {
	ret := golang.NewBlock("extendMutation", "graphql")
	ret.W("extend type Mutation {")
	ret.W("}")
	return ret
}
