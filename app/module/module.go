package module

import (
	"fmt"
	"os"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

type Module struct {
	Key         string `json:"-"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	License     string `json:"license"`
	Sourcecode  string `json:"sourcecode"`
	Description string `json:"description"`
}

func (m *Module) Path() string {
	return fmt.Sprintf("module/%s", m.Key)
}

func (m *Module) FileContent(files filesystem.FileLoader, path string) (os.FileMode, []byte, error) {
	stat := files.Stat(path)
	if stat == nil {
		return 0, nil, errors.Errorf("file [%s] not found", path)
	}
	b, err := files.ReadFile(path)
	if err != nil {
		return 0, nil, err
	}
	return stat.Mode(), b, nil
}


type Modules map[string]*Module

func (i Modules) Get(key string) *Module {
	for _, item := range i {
		if item.Key == key {
			return item
		}
	}
	return nil
}

const (
	KeyModuleBootstrap = "bootstrap"
	KeyModuleCore      = "core"
)

var Bootstrap = &Module{
	Key:         KeyModuleBootstrap,
	AuthorName:  "Kyle U",
	AuthorEmail: "kyle@kyleu.com",
	License:     "Proprietary",
	Sourcecode:  "https://github.com/kyleu/projectforge/tree/master/module/bootstrap",
	Description: "Used to bootstrap new applications, don't use this directly",
}

var AvailableModules = Modules{
	KeyModuleBootstrap: Bootstrap,
	KeyModuleCore: {
		Key:         KeyModuleCore,
		AuthorName:  "Kyle U",
		AuthorEmail: "kyle@kyleu.com",
		License:     "Proprietary",
		Sourcecode:  "https://github.com/kyleu/projectforge/tree/master/module/core",
		Description: "Core module from [" + util.AppName + "]",
	},
}
