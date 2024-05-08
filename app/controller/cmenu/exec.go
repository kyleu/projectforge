// Package cmenu - Content managed by Project Forge, see [projectforge.md] for details.
package cmenu

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/lib/menu"
)

func processMenu(processes exec.Execs) *menu.Item {
	ret := make(menu.Items, 0, len(processes))
	lo.ForEach(processes, func(p *exec.Exec, _ int) {
		title := p.String()
		if p.Completed != nil {
			title += "*"
		}
		ret = append(ret, &menu.Item{Key: p.String(), Title: title, Icon: "bolt", Description: p.String(), Route: p.WebPath()})
	})
	desc := "process executions managed by this system"
	return &menu.Item{Key: "exec", Title: "Processes", Description: desc, Icon: "desktop", Route: "/admin/exec", Children: ret}
}
