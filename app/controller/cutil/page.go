// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cmenu"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

const (
	DefaultSearchPath  = "/search"
	DefaultProfilePath = "/profile"
	defaultIcon        = "app"
)

var (
	defaultRootTitleAppend = util.GetEnv("app_display_name_append")
	defaultRootTitle       = func() string {
		if tmp := util.GetEnv("app_display_name"); tmp != "" {
			return tmp
		}
		return util.AppName
	}()
)

type PageState struct {
	Title         string            `json:"title,omitempty"`
	Description   string            `json:"description,omitempty"`
	Method        string            `json:"method,omitempty"`
	URI           *fasthttp.URI     `json:"-"`
	Menu          menu.Items        `json:"menu,omitempty"`
	Breadcrumbs   cmenu.Breadcrumbs `json:"breadcrumbs,omitempty"`
	Flashes       []string          `json:"flashes,omitempty"`
	Session       util.ValueMap     `json:"-"`
	Profile       *user.Profile     `json:"profile,omitempty"`
	Accounts      user.Accounts     `json:"accounts,omitempty"`
	Authed        bool              `json:"authed,omitempty"`
	Admin         bool              `json:"admin,omitempty"`
	Icons         []string          `json:"icons,omitempty"`
	RootIcon      string            `json:"rootIcon,omitempty"`
	RootPath      string            `json:"rootPath,omitempty"`
	RootTitle     string            `json:"rootTitle,omitempty"`
	SearchPath    string            `json:"searchPath,omitempty"`
	ProfilePath   string            `json:"profilePath,omitempty"`
	HideMenu      bool              `json:"hideMenu,omitempty"`
	ForceRedirect string            `json:"forceRedirect,omitempty"`
	HeaderContent string            `json:"headerContent,omitempty"`
	Data          any               `json:"data,omitempty"`
	Logger        util.Logger       `json:"-"`
	Context       context.Context   `json:"-"`
	Span          *telemetry.Span   `json:"-"`
	RenderElapsed float64           `json:"renderElapsed,omitempty"`
}

func (p *PageState) AddIcon(n string) {
	for _, icon := range p.Icons {
		if icon == n {
			return
		}
	}
	p.Icons = append(p.Icons, n)
}

func (p *PageState) TitleString() string {
	if p.Title == "" {
		return util.AppName
	}
	return fmt.Sprintf("%s - %s", p.Title, util.AppName)
}

func (p *PageState) User() string {
	if len(p.Accounts) == 0 {
		return "anonymous"
	}
	return p.Accounts[0].Email
}

func (ps *PageState) Clean(as *app.State) error {
	if ps.Profile != nil && ps.Profile.Theme == "" {
		ps.Profile.Theme = theme.ThemeDefault.Key
	}
	if ps.RootIcon == "" {
		ps.RootIcon = defaultIcon
	}
	if ps.RootPath == "" {
		ps.RootPath = "/"
	}
	if ps.RootTitle == "" {
		ps.RootTitle = defaultRootTitle
	}
	if defaultRootTitleAppend != "" {
		ps.RootTitle += " " + defaultRootTitleAppend
	}
	if ps.SearchPath == "" {
		ps.SearchPath = DefaultSearchPath
	}
	if ps.ProfilePath == "" {
		ps.ProfilePath = DefaultProfilePath
	}
	if len(ps.Menu) == 0 {
		m, err := cmenu.MenuFor(ps.Context, ps.Authed, ps.Admin, as, ps.Logger)
		if err != nil {
			return err
		}
		ps.Menu = m
	}
	return nil
}

func (p *PageState) Close() {
	if p.Span != nil {
		p.Span.Complete()
	}
}
