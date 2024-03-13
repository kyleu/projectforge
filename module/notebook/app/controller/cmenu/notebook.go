package cmenu

import (
	"context"

	"{{{ .Package }}}/app/lib/menu"
)

func notebookMenu(ctx context.Context, showFiles bool) *menu.Item {
	var kids menu.Items
	if showFiles {
		kids = append(kids, &menu.Item{Key: "files", Title: "Files", Description: "Notebook files", Icon: "file", Route: "/notebook/files"})
	}
	return &menu.Item{Key: "notebook", Title: "Notebook", Description: "A Observable Framework notebook", Icon: "notebook", Route: "/notebook", Children: kids}
}
