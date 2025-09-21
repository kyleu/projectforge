package project

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project/export/load"
	"projectforge.dev/projectforge/app/util"
)

var DefaultIcon = "code"

type Project struct {
	Key     string   `json:"key"`
	Name    string   `json:"name,omitzero"`
	Icon    string   `json:"icon,omitzero"`
	Exec    string   `json:"exec,omitzero"`
	Version string   `json:"version"`
	Package string   `json:"package,omitzero"`
	Args    string   `json:"args,omitzero"`
	Port    int      `json:"port,omitzero"`
	Modules []string `json:"modules"`
	Ignore  []string `json:"ignore,omitempty"`
	Tags    []string `json:"tags"`

	Info  *Info        `json:"info,omitzero"`
	Theme *theme.Theme `json:"theme,omitzero"`
	Build *Build       `json:"build,omitzero"`
	Files []string     `json:"files,omitempty"`

	ExportArgs *metamodel.Args `json:"-"`
	Config     util.ValueMap   `json:"-"`
	Path       string          `json:"-"`
	Parent     string          `json:"-"`
	Error      string          `json:"error,omitzero"`
}

func NewProject(key string, path string) *Project {
	_, key = util.StringSplitPath(key)
	return &Project{Key: key, Version: "0.0.0", Path: path, Modules: []string{"core"}, Theme: theme.Default}
}

func (p *Project) Title() string {
	return util.OrDefault(p.Name, p.Key)
}

func (p *Project) NameSafe() string {
	return strings.ReplaceAll(strings.ReplaceAll(p.Name, "-", ""), " ", "")
}

func (p *Project) Executable() string {
	return util.OrDefault(p.Exec, p.Key)
}

func (p *Project) CleanKey() string {
	return CleanKey(p.Key)
}

func (p *Project) ExecSafe() string {
	if p.Exec != "" {
		return p.Exec
	}
	return p.Key
}

func (p *Project) IconSafe() string {
	if _, ok := util.SVGLibrary[p.Icon]; ok {
		return p.Icon
	}
	return DefaultIcon
}

func (p *Project) DescriptionSafe() string {
	if p.Info == nil {
		return ""
	}
	return util.OrDefault(p.Info.Description, p.Info.Summary)
}

func (p *Project) HasModule(key string) bool {
	return lo.Contains(p.Modules, key)
}

func (p *Project) ToMap() util.ValueMap {
	return util.ValueMap{
		"key": p.Key, "name": p.Name, "icon": p.Icon, "exec": p.Exec,
		"version": p.Version, "package": p.Package, "args": p.Key, "port": p.Port,
		"modules": p.Modules, "path": p.Path,
	}
}

func (p *Project) WebPath() string {
	return "/p/" + p.Key
}

func (p *Project) WebPathEnums() string {
	return p.WebPath() + "/export/enums"
}

func (p *Project) WebPathEvents() string {
	return p.WebPath() + "/export/events"
}

func (p *Project) WebPathModels() string {
	return p.WebPath() + "/export/models"
}

func (p *Project) ModuleArgExport(pSvc *Service, logger util.Logger) error {
	if p.ExportArgs == nil && p.HasModule("export") {
		fs, err := pSvc.GetFilesystem(p)
		if err != nil {
			return err
		}
		p.ExportArgs, err = load.ExportArgs(fs, ConfigDir, p.Info.Acronyms, logger)
		if err != nil {
			return err
		}
		p.ExportArgs.Modules = p.Modules
		p.ExportArgs.Database = p.DatabaseEngineDefault()
		p.ExportArgs.ApplyAcronyms(p.Info.Acronyms...)
	}
	return nil
}

func (p *Project) GoVersion() string {
	if p.Info == nil || p.Info.GoVersion == "" {
		return DefaultGoVersion
	}
	return p.Info.GoVersion
}

var Fields = []string{"Key", "Name", "Icon", "Exec", "Version", "Package", "Args", "Port", "Modules", "Ignore", "Tags", "Theme", "Path", "Parent", "Error"}

func (p *Project) ToCSV() ([]string, [][]string) {
	return Fields, [][]string{p.Strings()}
}

func (p *Project) Strings() []string {
	return []string{
		p.Key, p.Name, p.Icon, p.Exec, p.Version, p.Package, p.Args, fmt.Sprint(p.Port), util.StringJoin(p.Modules, ","),
		util.StringJoin(p.Ignore, ","), util.StringJoin(p.Tags, ","), p.Theme.Key, p.Path, p.Parent, p.Error,
	}
}
