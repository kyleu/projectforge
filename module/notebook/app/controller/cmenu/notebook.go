package cmenu

import (
	"context"

	"{{{ .Package }}}/app/lib/menu"
)

func notebookMenu(_ context.Context) *menu.Item {
	return &menu.Item{Key: "notebook", Title: "Notebook", Description: "Notebook", Icon: "notebook", Route: "/notebook"}
}
