// Code generated by qtc from "Detail.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vtask/Detail.html:1
package vtask

//line views/vtask/Detail.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/task"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/components/edit"
	"projectforge.dev/projectforge/views/components/view"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vexec"
)

//line views/vtask/Detail.html:13
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vtask/Detail.html:13
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vtask/Detail.html:13
type Detail struct {
	layout.Basic
	Task      *task.Task
	Result    *task.Result
	Args      util.ValueMap
	SocketURL string
}

//line views/vtask/Detail.html:21
func (p *Detail) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vtask/Detail.html:21
	qw422016.N().S(`
  <div class="card">
    <h3>`)
//line views/vtask/Detail.html:23
	components.StreamSVGIcon(qw422016, p.Task.IconSafe(), ps)
//line views/vtask/Detail.html:23
	qw422016.N().S(` `)
//line views/vtask/Detail.html:23
	qw422016.E().S(p.Task.TitleSafe())
//line views/vtask/Detail.html:23
	qw422016.N().S(`</h3>
    <form action="`)
//line views/vtask/Detail.html:24
	qw422016.E().S(p.Task.WebPath())
//line views/vtask/Detail.html:24
	qw422016.N().S(`/run" method="get">
      <table class="mt expanded">
        <tbody>
          `)
//line views/vtask/Detail.html:27
	edit.StreamTableEditorNoTable(qw422016, p.Task.Key, p.Task.Fields, p.Args)
//line views/vtask/Detail.html:27
	qw422016.N().S(`
          <tr>
            <td colspan="2">
              <button type="submit">Run</button>
              <button type="submit" name="async" value="true">Start</button>
            </td>
          </tr>
        </tbody>
      </table>
    </form>
  </div>
`)
//line views/vtask/Detail.html:38
	if p.Result != nil {
//line views/vtask/Detail.html:38
		qw422016.N().S(`  `)
//line views/vtask/Detail.html:39
		StreamResult(qw422016, as, p.Result, ps)
//line views/vtask/Detail.html:39
		qw422016.N().S(`
`)
//line views/vtask/Detail.html:40
	}
//line views/vtask/Detail.html:40
	qw422016.N().S(`
`)
//line views/vtask/Detail.html:42
	if p.SocketURL != "" {
//line views/vtask/Detail.html:42
		qw422016.N().S(`  <script>
    function processMessage(m) {
      if (m.cmd === "complete") {
        const deets = document.getElementById("result-detail");
        deets.innerHTML = m.param.html;
      }
    }
  </script>
  `)
//line views/vtask/Detail.html:51
		StreamSocketContent(qw422016, as, "task-output", p.Task, p.SocketURL, ps, "processMessage(m)")
//line views/vtask/Detail.html:51
		qw422016.N().S(`
`)
//line views/vtask/Detail.html:52
	}
//line views/vtask/Detail.html:53
}

//line views/vtask/Detail.html:53
func (p *Detail) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vtask/Detail.html:53
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vtask/Detail.html:53
	p.StreamBody(qw422016, as, ps)
//line views/vtask/Detail.html:53
	qt422016.ReleaseWriter(qw422016)
//line views/vtask/Detail.html:53
}

//line views/vtask/Detail.html:53
func (p *Detail) Body(as *app.State, ps *cutil.PageState) string {
//line views/vtask/Detail.html:53
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vtask/Detail.html:53
	p.WriteBody(qb422016, as, ps)
//line views/vtask/Detail.html:53
	qs422016 := string(qb422016.B)
//line views/vtask/Detail.html:53
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vtask/Detail.html:53
	return qs422016
//line views/vtask/Detail.html:53
}

//line views/vtask/Detail.html:55
func StreamSocketContent(qw422016 *qt422016.Writer, as *app.State, key string, t *task.Task, u string, ps *cutil.PageState, callbacks ...string) {
//line views/vtask/Detail.html:55
	qw422016.N().S(`
  <div class="card">
    <div class="right">`)
//line views/vtask/Detail.html:57
	view.StreamTimestampRelative(qw422016, util.TimeCurrentP(), false)
//line views/vtask/Detail.html:57
	qw422016.N().S(`</div>
    <h3>`)
//line views/vtask/Detail.html:58
	components.StreamSVGIcon(qw422016, "file", ps)
//line views/vtask/Detail.html:58
	qw422016.N().S(` `)
//line views/vtask/Detail.html:58
	qw422016.E().S(t.TitleSafe())
//line views/vtask/Detail.html:58
	qw422016.N().S(` Logs</h3>
    <div class="mt">`)
//line views/vtask/Detail.html:59
	components.StreamTerminal(qw422016, "task-output", "Starting task ["+t.TitleSafe()+"]...")
//line views/vtask/Detail.html:59
	qw422016.N().S(`</div>
  </div>
  <div id="result-detail"></div>
  `)
//line views/vtask/Detail.html:62
	vexec.StreamExecScript(qw422016, key, u, as.Debug, ps, callbacks...)
//line views/vtask/Detail.html:62
	qw422016.N().S(`
`)
//line views/vtask/Detail.html:63
}

//line views/vtask/Detail.html:63
func WriteSocketContent(qq422016 qtio422016.Writer, as *app.State, key string, t *task.Task, u string, ps *cutil.PageState, callbacks ...string) {
//line views/vtask/Detail.html:63
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vtask/Detail.html:63
	StreamSocketContent(qw422016, as, key, t, u, ps, callbacks...)
//line views/vtask/Detail.html:63
	qt422016.ReleaseWriter(qw422016)
//line views/vtask/Detail.html:63
}

//line views/vtask/Detail.html:63
func SocketContent(as *app.State, key string, t *task.Task, u string, ps *cutil.PageState, callbacks ...string) string {
//line views/vtask/Detail.html:63
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vtask/Detail.html:63
	WriteSocketContent(qb422016, as, key, t, u, ps, callbacks...)
//line views/vtask/Detail.html:63
	qs422016 := string(qb422016.B)
//line views/vtask/Detail.html:63
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vtask/Detail.html:63
	return qs422016
//line views/vtask/Detail.html:63
}