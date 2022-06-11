package cmenu

import (
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/lib/menu"
)

func DoctorMenu(i string, r string) *menu.Item {
	return &menu.Item{Key: action.TypeDoctor.Key, Title: action.TypeDoctor.Title, Description: action.TypeDoctor.Description, Icon: i, Route: r}
}
