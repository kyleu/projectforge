package cmenu

import (
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/har"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
)

func harMenu(s *har.Service, logger util.Logger) *menu.Item {
	harKids := lo.Map(s.List(logger), func(n string, _ int) *menu.Item {
		n = strings.TrimSuffix(n, har.Ext)
		return &menu.Item{Key: n, Title: n, Icon: "book", Route: "/har/" + n}
	})
	return &menu.Item{Key: "har", Title: "Archives", Description: "HTTP Archive files", Icon: "book", Route: "/har", Children: harKids}
}
