package project

const (
	defaultType    = "projectforge.dev"
	currentVersion = "0.0.0"
)

type Project struct {
	Key     string   `json:"key"`
	Type    string   `json:"type"`
	Name    string   `json:"name,omitempty"`
	Version string   `json:"version"`
	Package string   `json:"package,omitempty"`
	Args    string   `json:"args,omitempty"`
	Port    int      `json:"port,omitempty"`
	Modules []string `json:"modules"`
	Ignore  []string `json:"ignore,omitempty"`
	Info    *Info    `json:"info"`
	Path    string   `json:"-"`
}

func NewProject(key string, path string) *Project {
	return &Project{Key: key, Type: defaultType, Version: currentVersion, Path: path}
}
