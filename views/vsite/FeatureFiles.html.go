// Code generated by qtc from "FeatureFiles.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vsite/FeatureFiles.html:1
package vsite

//line views/vsite/FeatureFiles.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vfile"
)

//line views/vsite/FeatureFiles.html:11
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vsite/FeatureFiles.html:11
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vsite/FeatureFiles.html:11
type FeatureFiles struct {
	layout.Basic
	Module *module.Module
	Path   []string
}

//line views/vsite/FeatureFiles.html:17
func (p *FeatureFiles) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vsite/FeatureFiles.html:17
	qw422016.N().S(`
`)
//line views/vsite/FeatureFiles.html:19
	mod := p.Module
	fs := as.Services.Modules.GetFilesystem(mod.Key)
	u := mod.FeaturesFilePath()

//line views/vsite/FeatureFiles.html:22
	qw422016.N().S(`
`)
//line views/vsite/FeatureFiles.html:24
	if fs.IsDir(util.StringFilePath(p.Path...)) {
//line views/vsite/FeatureFiles.html:25
		files := fs.ListFiles(util.StringFilePath(p.Path...), nil, ps.Logger)

//line views/vsite/FeatureFiles.html:25
		qw422016.N().S(`  <div class="card">
    <h3>`)
//line views/vsite/FeatureFiles.html:27
		components.StreamSVGIcon(qw422016, p.Module.Icon, ps)
//line views/vsite/FeatureFiles.html:27
		qw422016.N().S(` `)
//line views/vsite/FeatureFiles.html:27
		qw422016.E().S(p.Module.Title())
//line views/vsite/FeatureFiles.html:27
		qw422016.N().S(`</h3>
    `)
//line views/vsite/FeatureFiles.html:28
		vfile.StreamList(qw422016, p.Path, files, fs, u, as, ps)
//line views/vsite/FeatureFiles.html:28
		qw422016.N().S(`
  </div>
`)
//line views/vsite/FeatureFiles.html:30
	} else {
//line views/vsite/FeatureFiles.html:32
		b, err := fs.ReadFile(util.StringFilePath(p.Path...))
		if err != nil {
			panic(err)
		}

//line views/vsite/FeatureFiles.html:36
		qw422016.N().S(`  <div class="card">
    <h3>`)
//line views/vsite/FeatureFiles.html:38
		components.StreamSVGIcon(qw422016, p.Module.Icon, ps)
//line views/vsite/FeatureFiles.html:38
		qw422016.N().S(` `)
//line views/vsite/FeatureFiles.html:38
		qw422016.E().S(p.Module.Title())
//line views/vsite/FeatureFiles.html:38
		qw422016.N().S(`</h3>
    `)
//line views/vsite/FeatureFiles.html:39
		vfile.StreamDetail(qw422016, p.Path, b, u, nil, as, ps)
//line views/vsite/FeatureFiles.html:39
		qw422016.N().S(`
  </div>
`)
//line views/vsite/FeatureFiles.html:41
	}
//line views/vsite/FeatureFiles.html:42
}

//line views/vsite/FeatureFiles.html:42
func (p *FeatureFiles) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vsite/FeatureFiles.html:42
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vsite/FeatureFiles.html:42
	p.StreamBody(qw422016, as, ps)
//line views/vsite/FeatureFiles.html:42
	qt422016.ReleaseWriter(qw422016)
//line views/vsite/FeatureFiles.html:42
}

//line views/vsite/FeatureFiles.html:42
func (p *FeatureFiles) Body(as *app.State, ps *cutil.PageState) string {
//line views/vsite/FeatureFiles.html:42
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vsite/FeatureFiles.html:42
	p.WriteBody(qb422016, as, ps)
//line views/vsite/FeatureFiles.html:42
	qs422016 := string(qb422016.B)
//line views/vsite/FeatureFiles.html:42
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vsite/FeatureFiles.html:42
	return qs422016
//line views/vsite/FeatureFiles.html:42
}
