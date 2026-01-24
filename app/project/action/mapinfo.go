package action

import (
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func infoFromCfg(proto *project.Project, cfg util.ValueMap) *project.Info {
	str := func(key string, def string) string {
		return util.OrDefault(cfg.GetStringOpt(key), def)
	}

	i := proto.Info
	if i == nil {
		i = &project.Info{}
	}

	authorID := str("authorID", str("author_id", i.AuthorID))
	if authorID == "" {
		authorID = proto.Key
	}
	authorName := str("authorName", str("author_name", i.AuthorName))
	if authorName == "" {
		authorName = proto.Key
	}
	authorEmail := str("authorEmail", str("author_email", i.AuthorEmail))
	if authorEmail == "" {
		authorEmail = proto.Key + "@github.com"
	}

	homepage := str("homepage", i.Homepage)
	if homepage == "" {
		homepage = "https://github.com/" + proto.Key + "/" + proto.Key
	}
	license := str("license", i.License)
	if license == "" {
		license = "Proprietary"
	}

	md := i.ModuleDefs
	if x := cfg.GetStringOpt("moduleDefs"); x != "" {
		_ = util.FromJSON([]byte(x), &md)
	}

	cfgVars := i.ConfigVars
	if x := cfg.GetStringOpt("configVars"); x != "" {
		_ = util.FromJSON([]byte(x), &cfgVars)
	}

	deps := i.Dependencies
	if x := cfg.GetStringOpt("dependencies"); x != "" {
		_ = util.FromJSON([]byte(x), &deps)
	}

	docs := i.Docs
	if x := cfg.GetStringOpt("docs"); x != "" {
		_ = util.FromJSON([]byte(x), &docs)
	}

	additionalPorts := i.AdditionalPorts
	if x := cfg.GetStringOpt("additionalPorts"); x != "" {
		_ = util.FromJSON([]byte(x), &additionalPorts)
		if len(additionalPorts) == 0 {
			additionalPorts = nil
		}
	}

	return &project.Info{
		Org:             str("org", i.Org),
		AuthorID:        authorID,
		AuthorName:      authorName,
		AuthorEmail:     authorEmail,
		Team:            str("team", i.Team),
		License:         license,
		Homepage:        homepage,
		Sourcecode:      str("sourcecode", i.Sourcecode),
		Summary:         str("summary", i.Summary),
		Description:     str("description", i.Description),
		CI:              str("ci", i.CI),
		Homebrew:        str("homebrew", i.Homebrew),
		Bundle:          str("bundle", i.Bundle),
		SigningIdentity: str("signingIdentity", str("signing_identity", i.SigningIdentity)),
		DatabaseEngine:  str("databaseEngine", i.DatabaseEngine),
		Slack:           str("slack", i.Slack),
		Channels:        util.StringSplitAndTrim(str("channels", util.StringJoin(i.Channels, ", ")), ","),
		JavaPackage:     str("javaPackage", str("java_package", i.JavaPackage)),
		GoVersion:       str("goVersion", str("go_version", i.GoVersion)),
		GoBinary:        str("goBinary", str("go_binary", i.GoBinary)),
		ConfigVars:      cfgVars,
		AdditionalPorts: additionalPorts,
		ExtraFiles:      util.StringSplitAndTrim(str("extraFiles", str("extra_files", util.StringJoin(i.ExtraFiles, ", "))), ","),
		IgnoredFiles:    util.StringSplitAndTrim(str("ignoredFiles", str("ignored_files", util.StringJoin(i.IgnoredFiles, ", "))), ","),
		Deployments:     util.StringSplitAndTrim(str("deployments", util.StringJoin(i.Deployments, ", ")), ","),
		EnvVars:         util.StringSplitAndTrim(str("envvars", str("env_vars", util.StringJoin(i.EnvVars, ", "))), ","),
		DockerPackages:  util.StringSplitAndTrim(str("dockerPackages", str("docker_packages", util.StringJoin(i.DockerPackages, ", "))), ","),
		Dependencies:    deps,
		Docs:            docs,
		Acronyms:        util.StringSplitAndTrim(str("acronyms", util.StringJoin(i.Acronyms, ", ")), ","),
		ModuleDefs:      md,
	}
}
