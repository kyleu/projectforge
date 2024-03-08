package project

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

var DefaultIcon = "code"

type Project struct {
	Key     string   `json:"key"`
	Name    string   `json:"name,omitempty"`
	Icon    string   `json:"icon,omitempty"`
	Exec    string   `json:"exec,omitempty"`
	Version string   `json:"version"`
	Package string   `json:"package,omitempty"`
	Args    string   `json:"args,omitempty"`
	Port    int      `json:"port,omitempty"`
	Modules []string `json:"modules"`
	Ignore  []string `json:"ignore,omitempty"`
	Tags    []string `json:"tags"`

	Info  *Info        `json:"info,omitempty"`
	Theme *theme.Theme `json:"theme,omitempty"`
	Build *Build       `json:"build,omitempty"`

	ExportArgs *model.Args   `json:"-"`
	Config     util.ValueMap `json:"-"`
	Path       string        `json:"-"`
	Parent     string        `json:"-"`
	Error      string        `json:"error,omitempty"`
}

func NewProject(key string, path string) *Project {
	_, key = util.StringSplitPath(key)
	return &Project{Key: key, Version: "0.0.0", Path: path, Modules: []string{"core"}, Theme: theme.Default}
}

func (p *Project) Title() string {
	if p.Name == "" {
		return p.Key
	}
	return p.Name
}

func (p *Project) NameSafe() string {
	return strings.ReplaceAll(strings.ReplaceAll(p.Name, "-", ""), " ", "")
}

func (p *Project) Executable() string {
	if p.Exec == "" {
		return p.Key
	}
	return p.Exec
}

func (p *Project) CleanKey() string {
	return clean(p.Key)
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
	if p.Info.Description == "" {
		return p.Info.Summary
	}
	return p.Info.Description
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

func (p *Project) ModuleArgExport(pSvc *Service, logger util.Logger) (*model.Args, error) {
	if p.ExportArgs == nil && p.HasModule("export") {
		fs, err := pSvc.GetFilesystem(p)
		if err != nil {
			return nil, err
		}
		p.ExportArgs, err = pSvc.loadExportArgs(fs, logger)
		if err != nil {
			return nil, err
		}
		p.ExportArgs.Modules = p.Modules
		p.ExportArgs.Database = p.DatabaseEngineDefault()
		p.ExportArgs.Acronyms = p.Info.Acronyms
	}
	return p.ExportArgs, nil
}

func (p *Project) GoVersion() string {
	if p.Info == nil || p.Info.GoVersion == "" {
		return DefaultGoVersion
	}
	return p.Info.GoVersion
}

func (p *Project) DatabaseEngines() []string {
	ret := make([]string, 0, 4)
	if p.HasModule(util.DatabaseMySQL) {
		ret = append(ret, util.DatabaseMySQL)
	}
	if p.HasModule(util.DatabasePostgreSQL) {
		ret = append(ret, util.DatabasePostgreSQL)
	}
	if p.HasModule(util.DatabaseSQLite) {
		ret = append(ret, util.DatabaseSQLite)
	}
	if p.HasModule(util.DatabaseSQLServer) {
		ret = append(ret, util.DatabaseSQLServer)
	}
	return ret
}

func (p *Project) DatabaseEngineDefault() string {
	if p.Info != nil && p.Info.DatabaseEngine != "" {
		return p.Info.DatabaseEngine
	}
	engines := p.DatabaseEngines()
	if len(engines) == 1 {
		return engines[0]
	}
	if len(engines) > 1 {
		return util.DatabasePostgreSQL
	}
	return ""
}
