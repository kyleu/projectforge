package cutil

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cmenu"
	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/assets"
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
	Action         string            `json:"action,omitzero"`
	Title          string            `json:"title,omitzero"`
	Description    string            `json:"description,omitzero"`
	Method         string            `json:"method,omitzero"`
	URI            *url.URL          `json:"-"`
	Menu           menu.Items        `json:"menu,omitempty"`
	Breadcrumbs    cmenu.Breadcrumbs `json:"breadcrumbs,omitempty"`
	Flashes        []string          `json:"flashes,omitempty"`
	Session        util.ValueMap     `json:"-"`
	Profile        *user.Profile     `json:"profile,omitzero"`
	Authed         bool              `json:"authed,omitzero"`
	Admin          bool              `json:"admin,omitzero"`
	Params         filter.ParamSet   `json:"params,omitempty"`
	Icons          []string          `json:"icons,omitempty"`
	DefaultNavIcon string            `json:"defaultNavIcon,omitzero"`
	RootIcon       string            `json:"rootIcon,omitzero"`
	RootPath       string            `json:"rootPath,omitzero"`
	RootTitle      string            `json:"rootTitle,omitzero"`
	SearchPath     string            `json:"searchPath,omitzero"`
	ProfilePath    string            `json:"profilePath,omitzero"`
	HideHeader     bool              `json:"hideHeader,omitzero"`
	HideMenu       bool              `json:"hideMenu,omitzero"`
	NoStyle        bool              `json:"noStyle,omitzero"`
	NoScript       bool              `json:"noScript,omitzero"`
	ForceRedirect  string            `json:"forceRedirect,omitzero"`
	DefaultFormat  string            `json:"defaultFormat,omitzero"`
	HeaderContent  string            `json:"headerContent,omitzero"`
	Browser        string            `json:"browser,omitzero"`
	BrowserVersion string            `json:"browserVersion,omitzero"`
	OS             string            `json:"os,omitzero"`
	OSVersion      string            `json:"osVersion,omitzero"`
	Platform       string            `json:"platform,omitzero"`
	Transport      string            `json:"transport,omitzero"`
	Data           any               `json:"data,omitzero"`
	Started        time.Time         `json:"started,omitzero"`
	RenderElapsed  float64           `json:"renderElapsed,omitzero"`
	ResponseBytes  int64             `json:"responseBytes,omitzero"`
	ExtraContent   map[string]string `json:"extraContent,omitzero"`
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
}

func (p *PageState) Username() string {
	return p.Profile.Name
}

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
	}
	if p.SearchPath == "" {
		p.SearchPath = DefaultSearchPath
	}
	if p.ProfilePath == "" {
		p.ProfilePath = DefaultProfilePath
	}
	if len(p.Menu) == 0 {
		m, data, err := cmenu.MenuFor(p.Context, as, p.Authed, p.Admin, p.Profile, p.Params, p.Logger)
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
	classes := ret.JoinSpace()
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

func (p *PageState) Extra(key string) string {
	ret, ok := p.ExtraContent[key]
	return util.Choose(ok, ret, "")
}

func (p *PageState) AddHeaderScript(path string, deferFlag bool) {
	p.HeaderContent += "\n  " + assets.ScriptElement(path, deferFlag)
}

func (p *PageState) AddHeaderStylesheet(path string) {
	p.HeaderContent += "\n  " + assets.StylesheetElement(path)
}
