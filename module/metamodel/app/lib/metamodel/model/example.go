package model

import (
	"{{{ .Package }}}/app/lib/types"
)

var (
	ExampleColumns = Columns{
		{Name: "id", Type: types.NewUUID(), PK: true, Search: true, HelpString: "The primary key"},
		{Name: "name", Type: types.NewString(), Search: true, HelpString: "The name of the thing"},
		{Name: "created", Type: types.NewTimestamp(), SQLDefault: "now()", Tags: []string{"created"}, HelpString: "Created timestamp"},
		{Name: "deleted_at", Type: types.NewTimestamp(), Nullable: true, Tags: []string{"deleted"}, HelpString: "Optional timestamp"},
	}
	ExampleRelations = Relations{
		{Name: "relation_a", Src: []string{"parent_id"}, Table: "parent", Tgt: []string{"id"}},
	}
	ExampleIndexes = Indexes{
		{Name: "example_idx", Decl: `"table_name" ("id", "created")`},
	}
	ExampleEvent = map[string]any{
		"columns": ExampleColumns,
	}
	ExampleModel = map[string]any{
		"columns":   ExampleColumns,
		"relations": ExampleRelations,
		"indexes":   ExampleIndexes,
	}
)
