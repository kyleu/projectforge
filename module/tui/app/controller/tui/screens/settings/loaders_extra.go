package settings

import (
	"{{{ .Package }}}/app/controller/cmenu"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/lib/menu"
)

func sitemapLines(ts *mvc.State, _ *mvc.PageState) ([]string, error) {
	items, _, err := cmenu.MenuFor(ts.Context, ts.App, false, true, nil, nil, ts.Logger)
	if err != nil {
		return nil, err
	}
	return menuLines(items, ""), nil
}

func menuLines(items menu.Items, prefix string) []string {
	ret := make([]string, 0, len(items))
	for _, i := range items.Visible() {
		ret = append(ret, prefix+i.Title+" -> "+i.Route)
		if len(i.Children) > 0 {
			ret = append(ret, menuLines(i.Children, prefix+"  ")...)
		}
	}
	return ret
}
