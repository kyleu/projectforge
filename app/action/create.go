package action

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func onCreate(cfg util.ValueMap, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) *Result {
	ret := newResult(cfg, logger)
	path, _ := cfg.GetString("path", true)
	if path == "" {
		return errorResult(errors.New("must provide [path] as an argument"), cfg, logger)
	}
	if wipe, _ := cfg.GetBool("wipe"); wipe {
		fs := filesystem.NewFileSystem(".", logger)
		if fs.Exists(path) {
			ret.AddLog("removing existing directory [%s]", path)
			_ = fs.RemoveRecursive(path)
		}
	}

	prj := projectFromCfg(project.NewProject("new", path), cfg)

	err := writeConfig(prj, pSvc)
	if err != nil {
		return ret.WithError(err)
	}

	res, err := slam(prj, module.Bootstrap, false, mSvc, pSvc)
	if err != nil {
		return ret.WithError(err)
	}
	ret.Modules = append(ret.Modules, res)

	_, err = pSvc.Load(path)
	if err != nil {
		msg := fmt.Sprintf("unable to load newly created project from path [%s]", path)
		return errorResult(errors.Wrap(err, msg), cfg, logger)
	}

	prj, mods, err := prjAndMods(cfg, mSvc, pSvc)
	if err != nil {
		return errorResult(err, cfg, logger)
	}
	retS := onSlam(prj, mods, cfg, mSvc, pSvc, logger)
	ret = ret.Merge(retS)
	if ret.HasErrors() {
		return ret
	}

	ret = fullBuild(prj, ret, logger)

	return ret
}

func fullBuild(prj *project.Project, r *Result, logger *zap.SugaredLogger) *Result {
	r.AddLog("building project [%s] in [%s]", prj.Key, prj.Path)

	exitCode, out, err := util.RunProcessSimple("bin/templates.sh", prj.Path, logger)
	if err != nil {
		return r.WithError(err)
	}
	r.AddLog("templates.sh output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return r.WithError(errors.Errorf("templates.sh failed with exit code [%d]", exitCode))
	}

	exitCode, out, err = util.RunProcessSimple("go mod tidy", prj.Path, logger)
	if err != nil {
		return r.WithError(err)
	}
	r.AddLog("go mod output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return r.WithError(errors.Errorf("\"go mod tidy\" failed with exit code [%d]", exitCode))
	}

	exitCode, out, err = util.RunProcessSimple("npm install", filepath.Join(prj.Path, "client"), logger)
	if err != nil {
		return r.WithError(err)
	}
	r.AddLog("npm output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		r.WithError(errors.Errorf("npm install failed with exit code [%d]", exitCode))
	}

	exitCode, out, err = util.RunProcessSimple("bin/build/client.sh", prj.Path, logger)
	if err != nil {
		return r.WithError(err)
	}
	r.AddLog("client build output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		r.WithError(errors.Errorf("client build failed with exit code [%d]", exitCode))
	}

	exitCode, out, err = util.RunProcessSimple("make build", prj.Path, logger)
	if err != nil {
		return r.WithError(err)
	}
	r.AddLog("build output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		r.WithError(errors.Errorf("build failed with exit code [%d]", exitCode))
	}

	return r
}

func projectFromCfg(proto *project.Project, cfg util.ValueMap) *project.Project {
	str := func(key string, def ...string) string {
		ret, _ := cfg.GetString(key, true)
		if ret == "" && len(def) > 0 {
			ret = def[0]
		}
		return ret
	}
	integer := func(key string, def ...int) int {
		s := str(key)
		i, err := strconv.Atoi(s)
		if err != nil {
			return def[0]
		}
		return i
	}
	i := proto.Info
	if i == nil {
		i = &project.Info{License: "Proprietary"}
	}
	return &project.Project{
		Key:     str("key", proto.Key),
		Type:    proto.Type,
		Version: proto.Version,
		Name:    str("name", proto.Name),
		Package: str("package", proto.Package),
		Args:    str("args", proto.Args),
		Port:    integer("port", proto.Port),
		Modules: []string{"core"},
		Ignore:  proto.Ignore,
		Info: &project.Info{
			Org:         str("org", i.Org),
			AuthorName:  str("author_name", i.AuthorName),
			AuthorEmail: str("author_email", i.AuthorEmail),
			License:     str("license", i.License),
			Homepage:    str("homepage", i.Homepage),
			Sourcecode:  str("sourcecode", i.Sourcecode),
			Summary:     str("summary", i.Summary),
			Description: str("description", i.Description),
		},
		Path: proto.Path,
	}
}

func writeConfig(prj *project.Project, pSvc *project.Service) error {
	tgtFS := pSvc.GetFilesystem(prj)
	err := tgtFS.WriteFile(project.ConfigFilename, util.ToJSONBytes(prj, true), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write config file to [%s]", project.ConfigFilename)
	}
	return nil
}
