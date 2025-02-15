// Code generated by qtc from "Audit.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vaction/Audit.html:1
package vaction

//line views/vaction/Audit.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/views/components"
)

//line views/vaction/Audit.html:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vaction/Audit.html:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vaction/Audit.html:8
func StreamRenderAudit(qw422016 *qt422016.Writer, key string, res *action.AuditResult, as *app.State, ps *cutil.PageState) {
//line views/vaction/Audit.html:8
	qw422016.N().S(`
`)
//line views/vaction/Audit.html:9
	if res.Stats != nil {
//line views/vaction/Audit.html:9
		qw422016.N().S(`  `)
//line views/vaction/Audit.html:10
		StreamRenderCodeStats(qw422016, key, res.Stats, ps)
//line views/vaction/Audit.html:10
		qw422016.N().S(`
`)
//line views/vaction/Audit.html:11
	}
//line views/vaction/Audit.html:12
}

//line views/vaction/Audit.html:12
func WriteRenderAudit(qq422016 qtio422016.Writer, key string, res *action.AuditResult, as *app.State, ps *cutil.PageState) {
//line views/vaction/Audit.html:12
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vaction/Audit.html:12
	StreamRenderAudit(qw422016, key, res, as, ps)
//line views/vaction/Audit.html:12
	qt422016.ReleaseWriter(qw422016)
//line views/vaction/Audit.html:12
}

//line views/vaction/Audit.html:12
func RenderAudit(key string, res *action.AuditResult, as *app.State, ps *cutil.PageState) string {
//line views/vaction/Audit.html:12
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vaction/Audit.html:12
	WriteRenderAudit(qb422016, key, res, as, ps)
//line views/vaction/Audit.html:12
	qs422016 := string(qb422016.B)
//line views/vaction/Audit.html:12
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vaction/Audit.html:12
	return qs422016
//line views/vaction/Audit.html:12
}

//line views/vaction/Audit.html:14
func StreamRenderCodeStats(qw422016 *qt422016.Writer, key string, ci *action.CodeStats, ps *cutil.PageState) {
//line views/vaction/Audit.html:14
	qw422016.N().S(`
  <div class="overflow full-width mt">
    <div id="codestats_`)
//line views/vaction/Audit.html:16
	qw422016.E().S(key)
//line views/vaction/Audit.html:16
	qw422016.N().S(`"></div>
  </div>
  <div class="overflow full-width mt">
    <table class="expanded min-200">
      <thead>
        <tr>
          <th class="shrink">Type</th>
          <th>Code</th>
          <th>Comments</th>
          <th>Blanks</th>
          <th>Files</th>
        </tr>
      </thead>
      <tbody>
`)
//line views/vaction/Audit.html:30
	for _, t := range ci.Types {
//line views/vaction/Audit.html:30
		qw422016.N().S(`        <tr>
          <td>`)
//line views/vaction/Audit.html:32
		qw422016.E().S(t.Name)
//line views/vaction/Audit.html:32
		qw422016.N().S(`</td>
          <td>`)
//line views/vaction/Audit.html:33
		qw422016.N().D(t.Code)
//line views/vaction/Audit.html:33
		qw422016.N().S(`</td>
          <td>`)
//line views/vaction/Audit.html:34
		qw422016.N().D(t.Comments)
//line views/vaction/Audit.html:34
		qw422016.N().S(`</td>
          <td>`)
//line views/vaction/Audit.html:35
		qw422016.N().D(t.Blanks)
//line views/vaction/Audit.html:35
		qw422016.N().S(`</td>
          <td><a href="#modal-`)
//line views/vaction/Audit.html:36
		qw422016.E().S(key)
//line views/vaction/Audit.html:36
		qw422016.N().S(`-`)
//line views/vaction/Audit.html:36
		qw422016.E().S(t.Name)
//line views/vaction/Audit.html:36
		qw422016.N().S(`">`)
//line views/vaction/Audit.html:36
		qw422016.N().D(len(t.Files))
//line views/vaction/Audit.html:36
		qw422016.N().S(`</a></td>
        </tr>
`)
//line views/vaction/Audit.html:38
	}
//line views/vaction/Audit.html:39
	t := ci.Total

//line views/vaction/Audit.html:39
	qw422016.N().S(`        <tr style="border-top: var(--border);">
          <td style="border-top: var(--border);" class="nowrap"><em>`)
//line views/vaction/Audit.html:41
	qw422016.E().S(t.Name)
//line views/vaction/Audit.html:41
	qw422016.N().S(`</em></td>
          <td style="border-top: var(--border);">`)
//line views/vaction/Audit.html:42
	qw422016.N().D(t.Code)
//line views/vaction/Audit.html:42
	qw422016.N().S(`</td>
          <td style="border-top: var(--border);">`)
//line views/vaction/Audit.html:43
	qw422016.N().D(t.Comments)
//line views/vaction/Audit.html:43
	qw422016.N().S(`</td>
          <td style="border-top: var(--border);">`)
//line views/vaction/Audit.html:44
	qw422016.N().D(t.Blanks)
//line views/vaction/Audit.html:44
	qw422016.N().S(`</td>
          <td style="border-top: var(--border);"><a href="#modal-`)
//line views/vaction/Audit.html:45
	qw422016.E().S(key)
//line views/vaction/Audit.html:45
	qw422016.N().S(`-`)
//line views/vaction/Audit.html:45
	qw422016.E().S(t.Name)
//line views/vaction/Audit.html:45
	qw422016.N().S(`">`)
//line views/vaction/Audit.html:45
	qw422016.N().D(len(t.Files))
//line views/vaction/Audit.html:45
	qw422016.N().S(`</a></td>
        </tr>
      </tbody>
    </table>
  </div>
`)
//line views/vaction/Audit.html:50
	for _, t := range ci.Types {
//line views/vaction/Audit.html:50
		qw422016.N().S(`    `)
//line views/vaction/Audit.html:51
		streamcodeModal(qw422016, key, t, ps)
//line views/vaction/Audit.html:51
		qw422016.N().S(`
`)
//line views/vaction/Audit.html:52
	}
//line views/vaction/Audit.html:52
	qw422016.N().S(`  `)
//line views/vaction/Audit.html:53
	streamcodeModal(qw422016, key, ci.Total, ps)
//line views/vaction/Audit.html:53
	qw422016.N().S(`
  `)
//line views/vaction/Audit.html:54
	components.StreamPlotHorizontalBar(qw422016, "codestats_"+key, ci.ToMaps(), "files", "name", "["+key+"] Code Stats", 92)
//line views/vaction/Audit.html:54
	qw422016.N().S(`
`)
//line views/vaction/Audit.html:55
}

//line views/vaction/Audit.html:55
func WriteRenderCodeStats(qq422016 qtio422016.Writer, key string, ci *action.CodeStats, ps *cutil.PageState) {
//line views/vaction/Audit.html:55
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vaction/Audit.html:55
	StreamRenderCodeStats(qw422016, key, ci, ps)
//line views/vaction/Audit.html:55
	qt422016.ReleaseWriter(qw422016)
//line views/vaction/Audit.html:55
}

//line views/vaction/Audit.html:55
func RenderCodeStats(key string, ci *action.CodeStats, ps *cutil.PageState) string {
//line views/vaction/Audit.html:55
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vaction/Audit.html:55
	WriteRenderCodeStats(qb422016, key, ci, ps)
//line views/vaction/Audit.html:55
	qs422016 := string(qb422016.B)
//line views/vaction/Audit.html:55
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vaction/Audit.html:55
	return qs422016
//line views/vaction/Audit.html:55
}

//line views/vaction/Audit.html:57
func streamcodeModal(qw422016 *qt422016.Writer, key string, t *action.CodeType, ps *cutil.PageState) {
//line views/vaction/Audit.html:57
	qw422016.N().S(`
  <div id="modal-`)
//line views/vaction/Audit.html:58
	qw422016.E().S(key)
//line views/vaction/Audit.html:58
	qw422016.N().S(`-`)
//line views/vaction/Audit.html:58
	qw422016.E().S(t.Name)
//line views/vaction/Audit.html:58
	qw422016.N().S(`" class="modal" style="display: none;">
    <a class="backdrop" href="#"></a>
    <div class="modal-content" style="min-width: 90%;">
      <div class="modal-header">
        <a href="#" class="modal-close">×</a>
        <h2>`)
//line views/vaction/Audit.html:63
	qw422016.E().S(t.Name)
//line views/vaction/Audit.html:63
	qw422016.N().S(` Files</h2>
      </div>
      <div class="modal-body">
        <div class="overflow full-width">
          <table class="expanded min-200">
            <thead>
              <tr>
                <th>Name</th>
                <th>Code</th>
                <th>Comments</th>
                <th>Blanks</th>
                <th>Total</th>
              </tr>
            </thead>
            <tbody>
`)
//line views/vaction/Audit.html:78
	for _, x := range t.Files {
//line views/vaction/Audit.html:78
		qw422016.N().S(`            <tr>
              <td><a href="/p/`)
//line views/vaction/Audit.html:80
		qw422016.E().S(key)
//line views/vaction/Audit.html:80
		qw422016.N().S(`/fs/`)
//line views/vaction/Audit.html:80
		qw422016.E().S(x.Name)
//line views/vaction/Audit.html:80
		qw422016.N().S(`">`)
//line views/vaction/Audit.html:80
		qw422016.E().S(x.Name)
//line views/vaction/Audit.html:80
		qw422016.N().S(`</a></td>
              <td>`)
//line views/vaction/Audit.html:81
		qw422016.N().D(x.Code)
//line views/vaction/Audit.html:81
		qw422016.N().S(`</td>
              <td>`)
//line views/vaction/Audit.html:82
		qw422016.N().D(x.Comments)
//line views/vaction/Audit.html:82
		qw422016.N().S(`</td>
              <td>`)
//line views/vaction/Audit.html:83
		qw422016.N().D(x.Blanks)
//line views/vaction/Audit.html:83
		qw422016.N().S(`</td>
              <td>`)
//line views/vaction/Audit.html:84
		qw422016.N().D(x.Total())
//line views/vaction/Audit.html:84
		qw422016.N().S(`</td>
            </tr>
`)
//line views/vaction/Audit.html:86
	}
//line views/vaction/Audit.html:86
	qw422016.N().S(`            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
`)
//line views/vaction/Audit.html:93
}

//line views/vaction/Audit.html:93
func writecodeModal(qq422016 qtio422016.Writer, key string, t *action.CodeType, ps *cutil.PageState) {
//line views/vaction/Audit.html:93
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vaction/Audit.html:93
	streamcodeModal(qw422016, key, t, ps)
//line views/vaction/Audit.html:93
	qt422016.ReleaseWriter(qw422016)
//line views/vaction/Audit.html:93
}

//line views/vaction/Audit.html:93
func codeModal(key string, t *action.CodeType, ps *cutil.PageState) string {
//line views/vaction/Audit.html:93
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vaction/Audit.html:93
	writecodeModal(qb422016, key, t, ps)
//line views/vaction/Audit.html:93
	qs422016 := string(qb422016.B)
//line views/vaction/Audit.html:93
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vaction/Audit.html:93
	return qs422016
//line views/vaction/Audit.html:93
}
