package gql

import (
	"fmt"

	"github.com/samber/lo"
	"golang.org/x/exp/maps"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

type Module struct {
	Key         string `json:"key"`
	Name        string
	Icon        string
	Description string
	Hidden      bool
	AuthorName  string
	AuthorEmail string
	License     string
	Sourcecode  string
	ConfigVars  []string
	PortOffsets []string
	Dangerous   bool
	Requires    []string
	Priority    int32
	Technology  []string
	Files       []string
	URL         string
	UsageMD     []string
}

func FromModule(m *module.Module, logger util.Logger) *Module {
	files, _ := m.Files.ListFilesRecursive(".", nil, logger)
	ports := maps.Keys(lo.MapEntries(m.PortOffsets, func(k string, v int) (string, any) {
		return fmt.Sprintf("%s: %d", k, v), nil
	}))
	return &Module{
		Key:         m.Key,
		Name:        m.Name,
		Icon:        m.Icon,
		Description: m.Description,
		Hidden:      m.Hidden,
		AuthorName:  m.AuthorName,
		AuthorEmail: m.AuthorEmail,
		License:     m.License,
		Sourcecode:  m.Sourcecode,
		ConfigVars:  m.ConfigVars.Strings(),
		PortOffsets: ports,
		Dangerous:   m.Dangerous,
		Requires:    m.Requires,
		Priority:    int32(m.Priority),
		Technology:  m.Technology,
		Files:       files,
		URL:         m.URL,
		UsageMD:     util.StringSplitLines(m.UsageMD),
	}
}
