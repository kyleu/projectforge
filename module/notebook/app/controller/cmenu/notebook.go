package cmenu

import (
	"context"

	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/lib/notebook"
)

func notebookMenu(ctx context.Context, showFiles bool) *menu.Item {
	var kids menu.Items
	if showFiles {
		kids = append(kids, &menu.Item{Key: "files", Title: "Files", Description: "Notebook files", Icon: "folder", Route: "/notebook/files"})
		for _, page := range notebook.FavoritePages.Order {
			kids = append(kids, &menu.Item{Key: page, Title: notebook.FavoritePages.GetSimple(page), Icon: "file", Route: "/notebook/view/" + page})
		}
	}
	return &menu.Item{Key: "notebook", Title: "Notebook", Description: "A Observable Framework notebook", Icon: "notebook", Route: "/notebook", Children: kids}
}
