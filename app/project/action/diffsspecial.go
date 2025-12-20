package action

import (
	"fmt"
	"slices"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
)

const vscodeSettingsPath = ".vscode/settings.json"

func parseSettings(fs filesystem.FileLoader, thm *theme.Theme) (*file.File, error) {
	ret := util.ValueMap{}
	if fs.Exists(vscodeSettingsPath) {
		b, err := fs.ReadFile(vscodeSettingsPath)
		if err != nil {
			return nil, err
		}
		err = util.FromJSON(b, &ret)
		if err != nil {
			return nil, err
		}
	}
	colors := ret.GetMapOpt("workbench.colorCustomization")
	if colors == nil {
		colors = util.ValueMap{}
	}
	colors["titleBar.activeForeground"] = "#000000"
	colors["titleBar.inactiveForeground"] = "#000000aa"
	colors["titleBar.activeBackground"] = thm.Light.NavBackground
	colors["titleBar.inactiveBackground"] = thm.Light.MenuBackground
	ret["workbench.colorCustomizations"] = colors

	f := file.NewFile(vscodeSettingsPath, filesystem.DefaultMode, util.ToJSONBytes(ret, true))
	return f, nil
}

const vscodeLaunchPath = ".vscode/launch.json"

func parseLaunch(fs filesystem.FileLoader, name string, ex string, upgrade bool) (*file.File, error) {
	ret := util.ValueMap{}
	if fs.Exists(vscodeLaunchPath) {
		b, err := fs.ReadFile(vscodeLaunchPath)
		if err != nil {
			return nil, err
		}
		err = util.FromJSON(b, &ret)
		if err != nil {
			return nil, err
		}
	}
	config, _ := ret.GetMapArray("configurations", true)
	keys := lo.Map(config, func(x util.ValueMap, _ int) string {
		v := x["name"]
		if s, ok := v.(string); ok {
			return s
		}
		return fmt.Sprint(v)
	})
	add := func(l util.ValueMap) {
		if !slices.Contains(keys, l.GetStringOpt("name")) {
			config = append(config, l)
		}
	}
	sub := []any{util.ValueMap{"from": "${workspaceFolder}/views", "to": "views"}}
	add(util.ValueMap{"name": "Attach to " + name, "type": "go", "request": "attach", "mode": "local", "processId": ex, "substitutePath": sub})
	add(util.ValueMap{"name": "Start " + name, "type": "go", "request": "launch", "mode": "auto", "program": name, "substitutePath": sub})
	if upgrade {
		add(util.ValueMap{
			"name": "Start " + name, "type": "go", "request": "launch", "mode": "auto", "program": name, "args": []string{"upgrade"}, "substitutePath": sub,
		})
	}
	ret["configurations"] = config
	if ret.GetStringOpt("version") == "" {
		ret["version"] = "0.2.0"
	}
	f := file.NewFile(vscodeLaunchPath, filesystem.DefaultMode, util.ToJSONBytes(ret, true))
	return f, nil
}

func parseConfig(pm *PrjAndMods) (util.KeyTypeDescs, map[string]int) {
	var configVars util.KeyTypeDescs
	portOffsets := map[string]int{}

	lo.ForEach(pm.Prj.Modules, func(m string, _ int) {
		mod := pm.Mods.Get(m)
		lo.ForEach(mod.ConfigVars, func(src *util.KeyTypeDesc, _ int) {
			hit := lo.ContainsBy(configVars, func(tgt *util.KeyTypeDesc) bool {
				return src.Matches(tgt)
			})
			if !hit {
				configVars = append(configVars, src)
			}
		})
		for k, v := range mod.PortOffsets {
			portOffsets[k] = v
		}
	})
	for k, v := range pm.Prj.Info.AdditionalPorts {
		portOffsets[k] = v - pm.Prj.Port
	}
	return configVars.Sort(), portOffsets
}
