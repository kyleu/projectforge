package action

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

func diffs(prj *project.Project, mods module.Modules, addHeader bool, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) (file.Files, []*diff.Diff, error) {
	cs := toChangeset(prj)
	tgt := pSvc.GetFilesystem(prj)

	srcFiles, err := mSvc.GetFiles(mods, cs, addHeader, tgt)
	if err != nil {
		return nil, nil, err
	}

	diffs := diff.FileLoader(srcFiles, tgt, false, logger)

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
	add("NAME", p.Title())
	add("VERSION", p.Version)
	add("PACKAGE", p.Package)
	add("IGNORE_FILES", ignore)
	add("IGNORE_FILES_GREP", ignoreGrep)
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

	b := p.Build
	if p.Build == nil {
		b = &project.Build{}
	}
	addB := func(key string, tf bool) {
		add("BUILD_"+key, fmt.Sprintf("%t", tf))
	}
	addB("SKIP_DESKTOP", b.SkipDesktop)
	addB("SKIP_NOTARIZE", b.SkipNotarize)
	addB("SKIP_HOMEBREW", b.SkipHomebrew)

	addB("SKIP_WASM", b.SkipWASM)
	addB("SKIP_IOS", b.SkipIOS)
	addB("SKIP_ANDROID", b.SkipAndroid)

	addB("SKIP_LINUX_ARM", b.SkipLinuxArm)
	addB("SKIP_LINUX_MIPS", b.SkipLinuxMips)
	addB("SKIP_LINUX_ODD", b.SkipLinuxOdd)

	addB("SKIP_AIX", b.SkipAIX)
	addB("SKIP_DRAGONFLY", b.SkipDragonfly)
	addB("SKIP_ILLUMOS", b.SkipIllumos)
	addB("SKIP_FREEBSD", b.SkipFreeBSD)
	addB("SKIP_NETBSD", b.SkipNetBSD)
	addB("SKIP_OPENBSD", b.SkipOpenBSD)
	addB("SKIP_PLAN9", b.SkipPlan9)
	addB("SKIP_SOLARIS", b.SkipSolaris)

	rplc.Sort()
	return &file.Changeset{Replacements: rplc}
}
