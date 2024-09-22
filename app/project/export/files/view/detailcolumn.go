package view

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const stringLead = "{%%s "

func viewDetailColumn(
	g *golang.Template, ret *golang.Block, models model.Models, m *model.Model, link bool, col *model.Column, modelKey string, indent int, enums enum.Enums,
) {
	rels := m.RelationsFor(col)
	toStrings := strings.Join(viewDetailColumnString(rels, models, m, col), "")
	viewColumn(util.KeyDetail, g, ret, m, col, toStrings, link, modelKey, indent, models, enums, "p.Paths")
}

func ModelLinkURL(m *model.Model, prefix string, enums enum.Enums) string {
	return stringLead + prefix + "WebPath(paths...) %%}"
}

func viewDetailColumnString(rels model.Relations, models model.Models, m *model.Model, col *model.Column) []string {
	return lo.FilterMap(rels, func(rel *model.Relation, _ int) (string, bool) {
		relModel := models.Get(rel.Table)
		lCols := rel.SrcColumns(m)
		lNames := strings.Join(lCols.ProperNames(), "")

		x := fmt.Sprintf("p.%sBy%s", relModel.Proper(), lNames)
		relTitles := relModel.Columns.WithTag("title")
		if len(relTitles) == 0 {
			if relPks := relModel.PKs(); slices.Equal(relPks.Names(), rel.Tgt) {
				return x + "||", true
			}
		}
		msg := "%s||{%%%% if %s != nil %%%%} ({%%%%s %s.TitleString() %%%%})" + helper.TextEndIfExtra
		return fmt.Sprintf(msg, x, x, x), true
	})
}
