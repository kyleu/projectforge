package cmenu

import (
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/project/action"
)

func DoctorMenu(i string, r string) *menu.Item {
	return &menu.Item{Key: action.TypeDoctor.Key, Title: action.TypeDoctor.Title, Description: action.TypeDoctor.Description, Icon: i, Route: r}
}
