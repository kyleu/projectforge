package cutil

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/samber/lo"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cmenu"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/lib/theme"
	"{{{ .Package }}}/app/lib/user"{{{ if .HasUser }}}
	dbuser "{{{ .Package }}}/app/user"{{{ end }}}
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/assets"
)

const ({{{ if .HasModule "search" }}}
	DefaultSearchPath  = "/search"{{{ end }}}
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
	Action         string            `json:"action,omitempty"`
	Title          string            `json:"title,omitempty"`
	Description    string            `json:"description,omitempty"`
	Method         string            `json:"method,omitempty"`
	URI            *url.URL          `json:"-"`
	Menu           menu.Items        `json:"menu,omitempty"`
	Breadcrumbs    cmenu.Breadcrumbs `json:"breadcrumbs,omitempty"`
	Flashes        []string          `json:"flashes,omitempty"`
	Session        util.ValueMap     `json:"-"`{{{ if .HasUser }}}
	User           *dbuser.User      `json:"user,omitempty"`{{{ end }}}
	Profile        *user.Profile     `json:"profile,omitempty"`{{{ if .HasAccount }}}
	Accounts       user.Accounts     `json:"accounts,omitempty"`{{{ end }}}
	Authed         bool              `json:"authed,omitempty"`
	Admin          bool              `json:"admin,omitempty"`
	Params         filter.ParamSet   `json:"params,omitempty"`
	Icons          []string          `json:"icons,omitempty"`
	DefaultNavIcon string            `json:"defaultNavIcon,omitempty"`
	RootIcon       string            `json:"rootIcon,omitempty"`
	RootPath       string            `json:"rootPath,omitempty"`
	RootTitle      string            `json:"rootTitle,omitempty"`{{{ if .HasModule "search" }}}
	SearchPath     string            `json:"searchPath,omitempty"`{{{ end }}}
	ProfilePath    string            `json:"profilePath,omitempty"`
	HideHeader     bool              `json:"hideHeader,omitempty"`
	HideMenu       bool              `json:"hideMenu,omitempty"`
	NoScript       bool              `json:"noScript,omitempty"`
	ForceRedirect  string            `json:"forceRedirect,omitempty"`
	DefaultFormat  string            `json:"defaultFormat,omitempty"`
	HeaderContent  string            `json:"headerContent,omitempty"`
	Browser        string            `json:"browser,omitempty"`
	BrowserVersion string            `json:"browserVersion,omitempty"`
	OS             string            `json:"os,omitempty"`
	OSVersion      string            `json:"osVersion,omitempty"`
	Platform       string            `json:"platform,omitempty"`
	Data           any               `json:"data,omitempty"`
	Started        time.Time         `json:"started,omitempty"`
	RenderElapsed  float64           `json:"renderElapsed,omitempty"`
	ResponseBytes  int64             `json:"responseBytes,omitempty"`
	RequestBody    []byte            `json:"-"`
	W              *WriteCounter     `json:"-"`
	Logger         util.Logger       `json:"-"`
	Context        context.Context   `json:"-"` //nolint:containedctx // properly closed, never directly used
	Span           *telemetry.Span   `json:"-"`
}

func (p *PageState) AddIcon(keys ...string) {
	lo.ForEach(keys, func(k string, _ int) {
		if !lo.Contains(p.Icons, k) {
			p.Icons = append(p.Icons, k)
		}
	})
}

func (p *PageState) TitleString() string {
	if p.Title == "" {
		return util.AppName
	}
	if strings.HasPrefix(p.Title, "!") {
		return p.Title[1:]
	}
	return fmt.Sprintf("%s - %s", p.Title, util.AppName)
}{{{ if .HasUser }}}

func (p *PageState) Username() string {
	if p.User != nil {
		return p.User.Name
	}
	return p.Profile.Name
}

func (p *PageState) AuthString() string {
	n := p.Profile.String()
	if p.User != nil {
		n = p.User.Name
	}
	msg := fmt.Sprintf("signed in as %s", n){{{ if .HasAccount }}}
	if len(p.Accounts) == 0 {
		if n == user.DefaultProfile.Name {
			return "click to sign in"
		}
		return msg
	}
	return fmt.Sprintf("%s using [%s]", msg, p.Accounts.TitleString()){{{ else }}}
	return msg{{{ end }}}
}{{{ else }}}

func (p *PageState) Username() string {
	return p.Profile.Name
}{{{ if .HasAccount }}}

func (p *PageState) AuthString() string {
	n := p.Profile.String()
	msg := fmt.Sprintf("signed in as %s", n)
	if len(p.Accounts) == 0 {
		if n == user.DefaultProfile.Name {
			return "click to sign in"
		}
		return msg
	}
	return fmt.Sprintf("%s using [%s]", msg, p.Accounts.TitleString())
}{{{ end }}}{{{ end }}}

func (p *PageState) Clean(_ *http.Request, as *app.State) error {
	if p.Profile != nil && p.Profile.Theme == "" {
		p.Profile.Theme = theme.Default.Key
	}
	if p.RootIcon == "" {
		p.RootIcon = defaultIcon
	}
	if p.RootPath == "" {
		p.RootPath = "/"
	}
	if p.RootTitle == "" {
		p.RootTitle = defaultRootTitle
	}
	if defaultRootTitleAppend != "" {
		p.RootTitle += " " + defaultRootTitleAppend
	}{{{ if .HasModule "search" }}}
	if p.SearchPath == "" {
		p.SearchPath = DefaultSearchPath
	}{{{ end }}}
	if p.ProfilePath == "" {
		p.ProfilePath = DefaultProfilePath
	}
	if len(p.Menu) == 0 {
		m, data, err := cmenu.MenuFor(p.Context, p.Authed, p.Admin, p.Profile, p.Params, as, p.Logger)
		if err != nil {
			return err
		}
		if data != nil && p.Data == nil {
			p.Data = data
		}
		p.Menu = m
	}
	return nil
}

func (p *PageState) Close() {
	if p.Span != nil {
		p.Span.Complete()
	}
}

func (p *PageState) LogError(msg string, args ...any) {
	p.Logger.Errorf(msg, args...)
}

func (p *PageState) ClassDecl() string {
	if len(p.Icons) == 0 {
		return "-"
	}
	ret := &util.StringSlice{}
	if p.Profile.Mode != "" {
		ret.Push(p.Profile.ModeClass())
	}
	if p.Browser != "" {
		ret.Push("browser-" + p.Browser)
	}
	if p.OS != "" {
		ret.Push("os-" + p.OS)
	}
	if p.Platform != "" {
		ret.Push("platform-" + p.Platform)
	}
	if ret.Empty() {
		return ""
	}
	classes := ret.Join(" ")
	return fmt.Sprintf(` class=%q`, classes)
}

func (p *PageState) SetTitleAndData(title string, data any) {
	p.Title = title
	p.Data = data
}

func (p *PageState) MainClasses() string {
	var ret []string
	if p.HideHeader {
		ret = append(ret, "noheader")
	}
	if p.HideMenu {
		ret = append(ret, "nomenu")
	}
	return util.StringJoin(ret, " ")
}

func (p *PageState) AddHeaderScript(path string, deferFlag bool) {
	p.HeaderContent += "\n  " + assets.ScriptElement(path, deferFlag)
}

func (p *PageState) AddHeaderStylesheet(path string) {
	p.HeaderContent += "\n  " + assets.StylesheetElement(path)
}
