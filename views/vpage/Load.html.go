// Code generated by qtc from "Load.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// <!-- Content managed by Project Forge, see [projectforge.md] for details. -->

//line views/vpage/Load.html:2
package vpage

//line views/vpage/Load.html:2
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vpage/Load.html:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vpage/Load.html:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vpage/Load.html:8
type Load struct {
	layout.Basic
	URL              string
	Title            string
	Message          string
	HideInstructions bool
}

//line views/vpage/Load.html:16
func (p *Load) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vpage/Load.html:16
	qw422016.N().S(`
`)
//line views/vpage/Load.html:18
	if p.Message == "" {
		p.Message = "Please wait as your request is processed..."
	}

//line views/vpage/Load.html:21
	qw422016.N().S(`  <div class="card">
    <h3>`)
//line views/vpage/Load.html:23
	qw422016.E().S(p.Title)
//line views/vpage/Load.html:23
	qw422016.N().S(`</h3>
    <p>`)
//line views/vpage/Load.html:24
	qw422016.E().S(p.Message)
//line views/vpage/Load.html:24
	qw422016.N().S(`</p>
`)
//line views/vpage/Load.html:25
	if !p.HideInstructions {
//line views/vpage/Load.html:25
		qw422016.N().S(`    <div class="mt"><em>Please avoid refreshing the browser or navigating away, your page is loading</em></div>
`)
//line views/vpage/Load.html:27
	}
//line views/vpage/Load.html:27
	qw422016.N().S(`  </div>
  <meta http-equiv="refresh" content="0; url=`)
//line views/vpage/Load.html:29
	qw422016.E().S(p.URL)
//line views/vpage/Load.html:29
	qw422016.N().S(`">
`)
//line views/vpage/Load.html:30
}

//line views/vpage/Load.html:30
func (p *Load) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vpage/Load.html:30
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vpage/Load.html:30
	p.StreamBody(qw422016, as, ps)
//line views/vpage/Load.html:30
	qt422016.ReleaseWriter(qw422016)
//line views/vpage/Load.html:30
}

//line views/vpage/Load.html:30
func (p *Load) Body(as *app.State, ps *cutil.PageState) string {
//line views/vpage/Load.html:30
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vpage/Load.html:30
	p.WriteBody(qb422016, as, ps)
//line views/vpage/Load.html:30
	qs422016 := string(qb422016.B)
//line views/vpage/Load.html:30
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vpage/Load.html:30
	return qs422016
//line views/vpage/Load.html:30
}
