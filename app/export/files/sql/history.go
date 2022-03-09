package sql

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
)

func sqlHistory(ret *golang.Block, m *model.Model) {
	if m.IsHistory() {
		ret.W("")
		ret.W("create table if not exists %q (", m.Name+"_history")
		ret.W("  \"id\" uuid,")
		pkRefs := make([]string, 0, len(m.PKs()))
		for _, pk := range m.PKs() {
			x := fmt.Sprintf("\"%s_%s\"", m.Name, pk.Name)
			pkRefs = append(pkRefs, x)
			ret.W("  %s %s,", x, pk.ToSQLType())
		}
		ret.W("  \"o\" jsonb not null,")
		ret.W("  \"n\" jsonb not null,")
		ret.W("  \"c\" jsonb not null,")
		ret.W("  \"created\" timestamp not null default now(),")
		ret.W("  foreign key (%s) references %q (%s),", strings.Join(pkRefs, ", "), m.Name, strings.Join(m.PKs().NamesQuoted(), ", "))
		ret.W("  primary key (\"id\")")
		ret.W(");")
	}
}
