// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app/lib/menu"
	"github.com/kyleu/projectforge/app/lib/telemetry"
	"github.com/kyleu/projectforge/app/lib/user"
	"github.com/kyleu/projectforge/app/util"
)

type PageState struct {
	Title         string             `json:"title,omitempty"`
	Description   string             `json:"description,omitempty"`
	Method        string             `json:"method,omitempty"`
	URI           *fasthttp.URI      `json:"-"`
	Menu          menu.Items         `json:"menu,omitempty"`
	Breadcrumbs   Breadcrumbs        `json:"breadcrumbs,omitempty"`
	Flashes       []string           `json:"flashes,omitempty"`
	Session       util.ValueMap      `json:"-"`
	Profile       *user.Profile      `json:"profile,omitempty"`
	Accounts      user.Accounts      `json:"accounts,omitempty"`
	Authed        bool               `json:"authed,omitempty"`
	Admin         bool               `json:"admin,omitempty"`
	Icons         []string           `json:"icons,omitempty"`
	RootIcon      string             `json:"rootIcon,omitempty"`
	RootPath      string             `json:"rootPath,omitempty"`
	RootTitle     string             `json:"rootTitle,omitempty"`
	SearchPath    string             `json:"searchPath,omitempty"`
	ProfilePath   string             `json:"profilePath,omitempty"`
	HideMenu      bool               `json:"hideMenu,omitempty"`
	ForceRedirect string             `json:"forceRedirect,omitempty"`
	Data          interface{}        `json:"data,omitempty"`
	Logger        *zap.SugaredLogger `json:"-"`
	Context       context.Context    `json:"-"`
	Span          *telemetry.Span    `json:"-"`
	RenderElapsed float64            `json:"renderElapsed,omitempty"`
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

func (p *PageState) Close() {
	if p.Span != nil {
		p.Span.Complete()
	}
}
