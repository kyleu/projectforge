package sql

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func sqlHistory(ret *golang.Block, m *model.Model, modules []string) {
	if m.IsHistory() {
		ret.WB()
		ret.W("create table if not exists %q (", m.Name+"_history")
		ret.W("  \"id\" uuid,")
		pkRefs := make([]string, 0, len(m.PKs()))
		lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
			x := fmt.Sprintf("\"%s_%s\"", m.Name, pk.Name)
			pkRefs = append(pkRefs, x)
			st, _ := pk.ToSQLType()
			ret.W("  %s %s,", x, st)
		})
		ret.W("  \"o\" jsonb not null,")
		ret.W("  \"n\" jsonb not null,")
		ret.W("  \"c\" jsonb not null,")
		now := "now()"
		if lo.Contains(modules, "sqlite") && !lo.Contains(modules, "postgres") {
			now = "current_timestamp"
		}
		ret.W("  \"created\" timestamp not null default %s,", now)
		ret.W("  foreign key (%s) references %q (%s),", strings.Join(pkRefs, ", "), m.Name, strings.Join(m.PKs().NamesQuoted(), ", "))
		ret.W("  primary key (\"id\")")
		ret.W(");")
	}
}
