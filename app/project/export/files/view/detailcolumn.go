package view

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const stringLead = "{%%s "

func viewDetailColumn(ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, indent int, enums enum.Enums) {
	rels := m.RelationsFor(col)
	call, err := viewDetailColumnString(rels, models, m)
	if err != nil {
		panic(err)
	}
	viewColumn(util.KeyDetail, ret, m, col, call, link, modelKey, indent, models, enums, false, "p.Paths")
}

func ModelLinkURL(prefix string) string {
	return stringLead + prefix + "WebPath(paths...) %%}"
}

func viewDetailColumnString(rels model.Relations, models model.Models, m *model.Model) (string, error) {
	if len(rels) == 0 {
		return "", nil
	}
	if len(rels) > 1 {
		return "", errors.Errorf("expected one relation, found [%d]", len(rels))
	}
	rel := rels[0]
	relModel := models.Get(rel.Table)
	lCols := rel.SrcColumns(m)
	lNames := util.StringJoin(lCols.ProperNames(), "")

	return fmt.Sprintf("p.%sBy%s", relModel.Proper(), lNames), nil
}
