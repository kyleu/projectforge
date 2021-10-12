// Code generated by qtc from "List.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vproject/List.html:1
package vproject

//line views/vproject/List.html:1
import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/views/layout"
)

//line views/vproject/List.html:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vproject/List.html:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vproject/List.html:8
type List struct {
	layout.Basic
	Projects project.Projects
}

//line views/vproject/List.html:13
func (p *List) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vproject/List.html:13
	qw422016.N().S(`
  `)
//line views/vproject/List.html:14
	StreamTable(qw422016, p.Projects, as, ps)
//line views/vproject/List.html:14
	qw422016.N().S(`

  <div class="card">
    <h3>Actions for all Projects</h3>
    <p><a href="/run/preview"><button>Preview All</button></a></p>
  </div>
`)
//line views/vproject/List.html:20
}

//line views/vproject/List.html:20
func (p *List) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vproject/List.html:20
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vproject/List.html:20
	p.StreamBody(qw422016, as, ps)
//line views/vproject/List.html:20
	qt422016.ReleaseWriter(qw422016)
//line views/vproject/List.html:20
}

//line views/vproject/List.html:20
func (p *List) Body(as *app.State, ps *cutil.PageState) string {
//line views/vproject/List.html:20
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vproject/List.html:20
	p.WriteBody(qb422016, as, ps)
//line views/vproject/List.html:20
	qs422016 := string(qb422016.B)
//line views/vproject/List.html:20
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vproject/List.html:20
	return qs422016
//line views/vproject/List.html:20
}