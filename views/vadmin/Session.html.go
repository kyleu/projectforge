// Code generated by qtc from "Session.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// <!-- Content managed by Project Forge, see [projectforge.md] for details. -->

//line views/vadmin/Session.html:2
package vadmin

//line views/vadmin/Session.html:2
import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vadmin/Session.html:14
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vadmin/Session.html:14
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vadmin/Session.html:14
type Session struct{ layout.Basic }

//line views/vadmin/Session.html:16
func (p *Session) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vadmin/Session.html:16
	qw422016.N().S(`
  <div class="card">
    <h3>Session</h3>
    <em>`)
//line views/vadmin/Session.html:19
	qw422016.N().D(len(ps.Session))
//line views/vadmin/Session.html:19
	qw422016.N().S(` values</em>
  </div>
`)
//line views/vadmin/Session.html:21
	if len(ps.Session) > 0 {
//line views/vadmin/Session.html:21
		qw422016.N().S(`  <div class="card">
    <h3>Values</h3>
    <div class="overflow full-width">
      <table class="mt expanded">
        <tbody>
`)
//line views/vadmin/Session.html:27
		for _, k := range util.ArraySorted(lo.Keys(ps.Session)) {
//line views/vadmin/Session.html:28
			v := ps.Session[k]

//line views/vadmin/Session.html:28
			qw422016.N().S(`            <tr>
              <th class="shrink">`)
//line views/vadmin/Session.html:30
			qw422016.E().S(k)
//line views/vadmin/Session.html:30
			qw422016.N().S(`</th>
              <td>`)
//line views/vadmin/Session.html:31
			qw422016.E().S(fmt.Sprint(v))
//line views/vadmin/Session.html:31
			qw422016.N().S(`</td>
            </tr>
`)
//line views/vadmin/Session.html:33
		}
//line views/vadmin/Session.html:33
		qw422016.N().S(`        </tbody>
      </table>
    </div>
  </div>
`)
//line views/vadmin/Session.html:38
	} else {
//line views/vadmin/Session.html:38
		qw422016.N().S(`  <div class="card">
    <em>Empty session</em>
  </div>
`)
//line views/vadmin/Session.html:42
	}
//line views/vadmin/Session.html:42
	qw422016.N().S(`  <div class="card">
    <h3>Profile</h3>
    <div class="mt">`)
//line views/vadmin/Session.html:45
	components.StreamJSON(qw422016, ps.Profile)
//line views/vadmin/Session.html:45
	qw422016.N().S(`</div>
  </div>
`)
//line views/vadmin/Session.html:47
}

//line views/vadmin/Session.html:47
func (p *Session) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vadmin/Session.html:47
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vadmin/Session.html:47
	p.StreamBody(qw422016, as, ps)
//line views/vadmin/Session.html:47
	qt422016.ReleaseWriter(qw422016)
//line views/vadmin/Session.html:47
}

//line views/vadmin/Session.html:47
func (p *Session) Body(as *app.State, ps *cutil.PageState) string {
//line views/vadmin/Session.html:47
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vadmin/Session.html:47
	p.WriteBody(qb422016, as, ps)
//line views/vadmin/Session.html:47
	qs422016 := string(qb422016.B)
//line views/vadmin/Session.html:47
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vadmin/Session.html:47
	return qs422016
//line views/vadmin/Session.html:47
}
