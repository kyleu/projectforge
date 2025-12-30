package ctest

import (
	"fmt"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/layout"
)

func readmesTest(as *app.State, ps *cutil.PageState) (layout.Page, error) {
	persist := cutil.QueryStringBool(ps.URI, "persist")
	mods := as.Services.Modules.ModulesSorted()
	ret := make([]string, 0, len(mods))
	for _, mod := range mods {
		rm := readmeFor(mod)
		if persist {
			modFS := as.Services.Modules.GetFilesystem(mod.Key)
			err := modFS.WriteFile("README.md", []byte(rm), filesystem.DefaultMode, true)
			if err != nil {
				return nil, err
			}
		}
		ret = append(ret, rm)
	}
	ps.SetTitleAndData("Module Readmes", ret)
	page := &views.Debug{}
	return page, nil
}

const readmeContent = `# [%[2]s] Project Forge Module

This directory contains the files used by the %[1]q module of [Project Forge](https://projectforge.dev).

## Purpose

%[3]s

## Usage

When the %[1]q module is enabled in your project, all of the files in this directory (except this readme) will be processed and included in your application.

See [doc/module/%[1]s.md](doc/module/%[1]s.md) for usage information.
`

func readmeFor(m *module.Module) string {
	return fmt.Sprintf(readmeContent, m.Key, m.Name, m.Description)
}
