// Code generated by qtc from "FeatureList.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vsite/FeatureList.html:1
package vsite

//line views/vsite/FeatureList.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vmodule"
)

//line views/vsite/FeatureList.html:11
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vsite/FeatureList.html:11
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vsite/FeatureList.html:11
type FeatureList struct {
	layout.Basic
	Modules module.Modules
}

//line views/vsite/FeatureList.html:16
func (p *FeatureList) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vsite/FeatureList.html:16
	qw422016.N().S(`
  <div class="card">
    <h3>`)
//line views/vsite/FeatureList.html:18
	components.StreamSVGIcon(qw422016, `app`, ps)
//line views/vsite/FeatureList.html:18
	qw422016.N().S(` `)
//line views/vsite/FeatureList.html:18
	qw422016.E().S(util.AppName)
//line views/vsite/FeatureList.html:18
	qw422016.N().S(` Features</h3>
  </div>
`)
//line views/vsite/FeatureList.html:20
	for _, mod := range p.Modules {
//line views/vsite/FeatureList.html:20
		qw422016.N().S(`    <a class="link-section" href="/features/`)
//line views/vsite/FeatureList.html:21
		qw422016.E().S(mod.Key)
//line views/vsite/FeatureList.html:21
		qw422016.N().S(`">
      <div class="card">
        `)
//line views/vsite/FeatureList.html:23
		vmodule.StreamModuleTechList(qw422016, mod, ps)
//line views/vsite/FeatureList.html:23
		qw422016.N().S(`
        <div class="left mrs">`)
//line views/vsite/FeatureList.html:24
		components.StreamSVGRef(qw422016, mod.IconSafe(), 40, 40, "", ps)
//line views/vsite/FeatureList.html:24
		qw422016.N().S(`</div>
        <strong class="highlight">`)
//line views/vsite/FeatureList.html:25
		qw422016.E().S(mod.Title())
//line views/vsite/FeatureList.html:25
		qw422016.N().S(`</strong>
        <div><em>`)
//line views/vsite/FeatureList.html:26
		qw422016.E().S(mod.Description)
//line views/vsite/FeatureList.html:26
		qw422016.N().S(`</em></div>
      </div>
    </a>
`)
//line views/vsite/FeatureList.html:29
	}
//line views/vsite/FeatureList.html:30
}

//line views/vsite/FeatureList.html:30
func (p *FeatureList) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vsite/FeatureList.html:30
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vsite/FeatureList.html:30
	p.StreamBody(qw422016, as, ps)
//line views/vsite/FeatureList.html:30
	qt422016.ReleaseWriter(qw422016)
//line views/vsite/FeatureList.html:30
}

//line views/vsite/FeatureList.html:30
func (p *FeatureList) Body(as *app.State, ps *cutil.PageState) string {
//line views/vsite/FeatureList.html:30
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vsite/FeatureList.html:30
	p.WriteBody(qb422016, as, ps)
//line views/vsite/FeatureList.html:30
	qs422016 := string(qb422016.B)
//line views/vsite/FeatureList.html:30
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vsite/FeatureList.html:30
	return qs422016
//line views/vsite/FeatureList.html:30
}
