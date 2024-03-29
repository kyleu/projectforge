// Code generated by qtc from "View.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vsvg/View.html:1
package vsvg

//line views/vsvg/View.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/svg"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vsvg/View.html:9
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vsvg/View.html:9
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vsvg/View.html:9
type View struct {
	layout.Basic
	Project *project.Project
	SVG     *svg.SVG
}

//line views/vsvg/View.html:15
func (p *View) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vsvg/View.html:15
	qw422016.N().S(`
  <div class="card">
    <h3>SVG Icon [`)
//line views/vsvg/View.html:17
	qw422016.E().S(p.SVG.Key)
//line views/vsvg/View.html:17
	qw422016.N().S(`]</h3>

    <div class="two-pane mt">
      <div class="l" style="max-width: 128px">`)
//line views/vsvg/View.html:20
	qw422016.N().S(p.SVG.Markup)
//line views/vsvg/View.html:20
	qw422016.N().S(`</div>
      <div class="r">
        <ul>
          <li><a href="/svg/`)
//line views/vsvg/View.html:23
	qw422016.E().S(p.Project.Key)
//line views/vsvg/View.html:23
	qw422016.N().S(`/`)
//line views/vsvg/View.html:23
	qw422016.E().S(p.SVG.Key)
//line views/vsvg/View.html:23
	qw422016.N().S(`/setapp" title="Overwrite the current app icon and favicon assets"><button>Set app icon</button></a></li>
          <li class="mt"><a class="link-confirm" href="/svg/`)
//line views/vsvg/View.html:24
	qw422016.E().S(p.Project.Key)
//line views/vsvg/View.html:24
	qw422016.N().S(`/`)
//line views/vsvg/View.html:24
	qw422016.E().S(p.SVG.Key)
//line views/vsvg/View.html:24
	qw422016.N().S(`/remove" title="Remove this icon from the application" data-message="Are you sure you want to remove this icon?"><button>Remove</button></a></li>
        </ul>
      </div>
    </div>
  </div>
  <div class="card">
    <h3>Source</h3>
`)
//line views/vsvg/View.html:31
	out, _ := cutil.FormatLang(p.SVG.Markup, "svg")

//line views/vsvg/View.html:31
	qw422016.N().S(`    `)
//line views/vsvg/View.html:32
	qw422016.N().S(out)
//line views/vsvg/View.html:32
	qw422016.N().S(`
  </div>
`)
//line views/vsvg/View.html:34
}

//line views/vsvg/View.html:34
func (p *View) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vsvg/View.html:34
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vsvg/View.html:34
	p.StreamBody(qw422016, as, ps)
//line views/vsvg/View.html:34
	qt422016.ReleaseWriter(qw422016)
//line views/vsvg/View.html:34
}

//line views/vsvg/View.html:34
func (p *View) Body(as *app.State, ps *cutil.PageState) string {
//line views/vsvg/View.html:34
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vsvg/View.html:34
	p.WriteBody(qb422016, as, ps)
//line views/vsvg/View.html:34
	qs422016 := string(qb422016.B)
//line views/vsvg/View.html:34
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vsvg/View.html:34
	return qs422016
//line views/vsvg/View.html:34
}
