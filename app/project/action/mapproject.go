package action

import (
	"strings"

	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func ProjectFromMap(prj *project.Project, m util.ValueMap, parseKey bool) error {
	clean := func(s string) string {
		if strings.Contains(s, "\\") {
			s = util.StringSplitLastOnly(s, '\\', true)
		} else if strings.Contains(s, "/") {
			s = util.StringSplitLastOnly(s, '/', true)
		}
		return s
	}
	get := func(k string, def string) util.Str {
		return m.GetRichStringOpt(k).OrDefault(def)
	}
	getSplit := func(k string, def []string) util.Strings {
		delim := "||"
		return get(k, util.StringJoin(def, delim)).SplitAndTrim(delim)
	}
	getInt := func(key string, def int) int {
		i, ok := get(key, "").ParseInt()
		return util.Choose(ok, i, def)
	}

	if parseKey {
		prj.Key = clean(string(get("key", prj.Name).OrDefault(prj.Key)))
	}

	prj.Name = clean(get("name", prj.Name).String())
	prj.Icon = get("icon", prj.Icon).String()
	prj.Exec = get("exec", prj.Exec).String()
	prj.Version = get("version", prj.Version).String()
	prj.Package = get("package", prj.Package).OrDefault("github.com/" + prj.Key + "/" + prj.Key).String()
	prj.Args = get("args", prj.Args).String()
	prj.Port = getInt("port", prj.Port)
	if prj.Port == 0 {
		prj.Port = 20000
	}
	prj.Modules = util.ArraySorted(getSplit("modules", prj.Modules).Strings())
	if len(prj.Modules) == 0 {
		prj.Modules = []string{"core"}
	}
	prj.Ignore = getSplit("ignore", prj.Ignore).Strings()
	prj.Tags = getSplit("tags", prj.Tags).Strings()

	prj.Info = infoFromCfg(prj, m)

	if _, ok := m["light-background"]; ok {
		prj.Theme = theme.ApplyMap(prj.Theme, m)
	}
	if prj.Theme.Equals(theme.Default) {
		prj.Theme = nil
	}

	prj.Build = project.BuildFromMap(m)
	if prj.Build.Empty() {
		prj.Build = nil
	}

	prj.Files = getSplit("files", prj.Files).Strings()
	prj.Path = get("path", prj.Path).String()

	return nil
}
