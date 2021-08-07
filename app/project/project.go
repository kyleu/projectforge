package project

type Project struct {
	Key      string   `json:"key"`
	Type     string   `json:"type"`
	Name     string   `json:"name,omitempty"`
	Exec     string   `json:"exec,omitempty"`
	Version  string   `json:"version"`
	Package  string   `json:"package,omitempty"`
	Args     string   `json:"args,omitempty"`
	Port     int      `json:"port,omitempty"`
	Modules  []string `json:"modules"`
	Ignore   []string `json:"ignore,omitempty"`
	Children []string `json:"children,omitempty"`
	Info     *Info    `json:"info,omitempty"`
	Build    *Build   `json:"build,omitempty"`
	Path     string   `json:"-"`
	Parent   string   `json:"-"`
}

func (p *Project) Title() string {
	if p.Name == "" {
		return p.Key
	}
	return p.Name
}

func (p *Project) Executable() string {
	if p.Exec == "" {
		return p.Key
	}
	return p.Exec
}

func NewProject(key string, path string) *Project {
	return &Project{Key: key, Type: "projectforge.dev", Version: "0.0.0", Path: path}
}

type Projects []*Project
