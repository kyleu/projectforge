package action

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func diffs(prj *project.Project, mod *module.Module, addHeader bool, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) (file.Files, []*file.Diff, error) {
	cs := toChangeset(prj)
	tgt := pSvc.GetFilesystem(prj)

	srcFiles, err := mSvc.GetFiles(mod, cs, addHeader, tgt)
	if err != nil {
		return nil, nil, err
	}

	diffs := srcFiles.DiffFileLoader(tgt, false, logger)

	return srcFiles, diffs, nil
}

func toChangeset(p *project.Project) *file.Changeset {
	port := "10000"
	if p.Port > 0 {
		port = fmt.Sprintf("%d", p.Port)
	}
	args := " -v --addr=0.0.0.0 all"
	if p.Args != "" {
		args = " " + p.Args
	}
	ignore := strings.Join(util.StringArrayQuoted(p.Ignore), ", ")
	if ignore != "" {
		ignore = ", " + ignore
	}

	var ignoreGrep string
	for _, ig := range p.Ignore {
		ignoreGrep += fmt.Sprintf(" | grep -v \\\\./%s", ig)
	}

	var rplc file.Replacements
	add := func(k string, v string) {
		rplc = append(rplc, &file.Replacement{K: k, V: v})
	}

	add("KEY", p.Key)
	add("EXECUTABLE", p.Key)
	add("NAME", p.Name)
	add("VERSION", p.Version)
	add("PACKAGE", p.Package)
	add("IGNORE", ignore)
	add("IGNORE_GREP", ignoreGrep)
	add("PORT", port)
	add("ARGS", args)

	if p.Info != nil {
		i := p.Info
		add("ORG", i.Org)
		add("AUTHOR_NAME", i.AuthorName)
		add("AUTHOR_EMAIL", i.AuthorEmail)
		add("LICENSE", i.License)
		add("BUNDLE", i.Bundle)
		add("SIGNING_IDENTITY", i.SigningIdentity)
		add("HOMEPAGE", i.Homepage)
		add("SOURCECODE", i.Sourcecode)
		add("SUMMARY", i.Summary)
		add("DESCRIPTION", i.Description)
	}
	rplc.Sort()
	return &file.Changeset{Replacements: rplc}
}
