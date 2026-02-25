package settings

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/controller/tui/layout"
	"{{{ .Package }}}/app/controller/tui/mvc"
	"{{{ .Package }}}/app/controller/tui/screens"
	"{{{ .Package }}}/app/controller/tui/style"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
)

type linkListScreen struct {
	key      string
	title    string
	load     func(*mvc.State) menu.Items
	openLink func(*mvc.State, *menu.Item) string
	items    menu.Items
}

func newModulesScreen() *linkListScreen {
	return &linkListScreen{
		key:   keyModules,
		title: "Go Modules",
		load: func(ts *mvc.State) menu.Items {
			items := menu.Items{
				{
					Key:         util.AppSource,
					Title:       fmt.Sprintf("%s v%s", util.AppName, ts.App.BuildInfo.Version),
					Description: "Primary module",
				},
			}
			for _, m := range util.DebugBuildInfo().Deps {
				items = append(items, &menu.Item{Key: m.Path, Title: m.Path, Description: m.Version})
			}
			return items
		},
		openLink: func(_ *mvc.State, item *menu.Item) string {
			return packagePathURL(item.Key)
		},
	}
}

func newRoutesScreen() *linkListScreen {
	return &linkListScreen{
		key:   keyRoutes,
		title: "HTTP Routes",
		load: func(ts *mvc.State) menu.Items {
			ret := menu.Items{}
			for _, method := range util.MapKeysSorted(cutil.AppRoutesList) {
				for _, pth := range cutil.AppRoutesList[method] {
					ret = append(ret, &menu.Item{Key: pth, Title: method + " " + pth})
				}
			}
			return ret
		},
		openLink: func(ts *mvc.State, item *menu.Item) string {
			return joinURLPath(ts.ServerURL, item.Key)
		},
	}
}

func (s *linkListScreen) Key() string { return s.key }

func (s *linkListScreen) Init(ts *mvc.State, ps *mvc.PageState) tea.Cmd {
	ps.Title = s.title
	s.items = s.load(ts)
	ps.Cursor = menuClamp(ps.Cursor, len(s.items))
	ps.SetStatusText(util.StringPlural(len(s.items), "Item"))
	return nil
}

func (s *linkListScreen) Update(ts *mvc.State, ps *mvc.PageState, msg tea.Msg) (mvc.Transition, tea.Cmd, error) {
	if delta, ok := menuDelta(msg); ok {
		ps.Cursor = menuClamp(ps.Cursor+delta, len(s.items))
		return mvc.Stay(), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && m.String() == "enter" && len(s.items) > 0 {
		item := s.items[menuClamp(ps.Cursor, len(s.items))]
		url := s.openLink(ts, item)
		if err := screens.OpenInBrowser(url); err != nil {
			return mvc.Stay(), nil, err
		}
		ps.SetStatus("Opened %s", url)
		return mvc.Stay(), nil, nil
	}
	if m, ok := msg.(tea.KeyMsg); ok && (m.String() == "esc" || m.String() == "backspace" || m.String() == "b") {
		return mvc.Pop(), nil, nil
	}
	return mvc.Stay(), nil, nil
}

func (s *linkListScreen) View(ts *mvc.State, ps *mvc.PageState, rects layout.Rects) string {
	st := style.New(ts.Theme)
	return renderPanel(st, ps.Title, renderMenuBody(s.items, ps.Cursor, st, rects), rects)
}

func (s *linkListScreen) Help(_ *mvc.State, _ *mvc.PageState) screens.HelpBindings {
	return screens.HelpBindings{Short: []string{"up/down: move", "enter: open in browser", "b: back"}}
}

func joinURLPath(base string, pth string) string {
	if base == "" {
		return pth
	}
	base = strings.TrimRight(base, "/")
	if strings.HasPrefix(pth, "http://") || strings.HasPrefix(pth, "https://") {
		return pth
	}
	if strings.HasPrefix(pth, "/") {
		return base + pth
	}
	return base + "/" + pth
}

func packagePathURL(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return "https://" + path
}
