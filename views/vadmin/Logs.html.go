// Code generated by qtc from "Logs.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vadmin/Logs.html:1
package vadmin

//line views/vadmin/Logs.html:1
import (
	"go.uber.org/zap/zapcore"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vadmin/Logs.html:11
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vadmin/Logs.html:11
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vadmin/Logs.html:11
type Logs struct {
	layout.Basic
	Logs []*zapcore.Entry
}

//line views/vadmin/Logs.html:16
func (p *Logs) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vadmin/Logs.html:16
	qw422016.N().S(`
  <div class="card">
    <h3>`)
//line views/vadmin/Logs.html:18
	components.StreamSVGIcon(qw422016, `filter`, ps)
//line views/vadmin/Logs.html:18
	qw422016.N().S(` Recent Logs</h3>
    `)
//line views/vadmin/Logs.html:19
	streamlogTable(qw422016, p.Logs)
//line views/vadmin/Logs.html:19
	qw422016.N().S(`
  </div>
`)
//line views/vadmin/Logs.html:21
}

//line views/vadmin/Logs.html:21
func (p *Logs) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vadmin/Logs.html:21
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vadmin/Logs.html:21
	p.StreamBody(qw422016, as, ps)
//line views/vadmin/Logs.html:21
	qt422016.ReleaseWriter(qw422016)
//line views/vadmin/Logs.html:21
}

//line views/vadmin/Logs.html:21
func (p *Logs) Body(as *app.State, ps *cutil.PageState) string {
//line views/vadmin/Logs.html:21
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vadmin/Logs.html:21
	p.WriteBody(qb422016, as, ps)
//line views/vadmin/Logs.html:21
	qs422016 := string(qb422016.B)
//line views/vadmin/Logs.html:21
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vadmin/Logs.html:21
	return qs422016
//line views/vadmin/Logs.html:21
}

//line views/vadmin/Logs.html:23
func streamlogTable(qw422016 *qt422016.Writer, logs []*zapcore.Entry) {
//line views/vadmin/Logs.html:23
	qw422016.N().S(`
  <div class="overflow full-width">
    <table class="mt">
      <thead>
        <tr>
          <th>Level</th>
          <th>Message</th>
          <th>Occurred</th>
        </tr>
      </thead>
      <tbody>
`)
//line views/vadmin/Logs.html:34
	for _, l := range logs {
//line views/vadmin/Logs.html:34
		qw422016.N().S(`        <tr>
          <td>
`)
//line views/vadmin/Logs.html:37
		lv := l.Level.String()

//line views/vadmin/Logs.html:38
		switch lv {
//line views/vadmin/Logs.html:39
		case "debug", "trace":
//line views/vadmin/Logs.html:39
			qw422016.N().S(`            <em>debug</em>
`)
//line views/vadmin/Logs.html:41
		case "error", "fatal":
//line views/vadmin/Logs.html:41
			qw422016.N().S(`            <div class="error">error</div>
`)
//line views/vadmin/Logs.html:43
		default:
//line views/vadmin/Logs.html:43
			qw422016.N().S(`            `)
//line views/vadmin/Logs.html:44
			qw422016.E().S(lv)
//line views/vadmin/Logs.html:44
			qw422016.N().S(`
`)
//line views/vadmin/Logs.html:45
		}
//line views/vadmin/Logs.html:45
		qw422016.N().S(`          </td>
          <td>`)
//line views/vadmin/Logs.html:47
		qw422016.E().S(l.Message)
//line views/vadmin/Logs.html:47
		qw422016.N().S(`</td>
          <td>`)
//line views/vadmin/Logs.html:48
		qw422016.E().S(util.TimeRelative(&l.Time))
//line views/vadmin/Logs.html:48
		qw422016.N().S(`</td>
        </tr>
`)
//line views/vadmin/Logs.html:50
	}
//line views/vadmin/Logs.html:50
	qw422016.N().S(`      </tbody>
    </table>
  </div>
`)
//line views/vadmin/Logs.html:54
}

//line views/vadmin/Logs.html:54
func writelogTable(qq422016 qtio422016.Writer, logs []*zapcore.Entry) {
//line views/vadmin/Logs.html:54
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vadmin/Logs.html:54
	streamlogTable(qw422016, logs)
//line views/vadmin/Logs.html:54
	qt422016.ReleaseWriter(qw422016)
//line views/vadmin/Logs.html:54
}

//line views/vadmin/Logs.html:54
func logTable(logs []*zapcore.Entry) string {
//line views/vadmin/Logs.html:54
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vadmin/Logs.html:54
	writelogTable(qb422016, logs)
//line views/vadmin/Logs.html:54
	qs422016 := string(qb422016.B)
//line views/vadmin/Logs.html:54
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vadmin/Logs.html:54
	return qs422016
//line views/vadmin/Logs.html:54
}
