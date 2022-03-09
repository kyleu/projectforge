// Code generated by qtc from "Detail.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vaction/Detail.html:1
package vaction

//line views/vaction/Detail.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
)

//line views/vaction/Detail.html:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vaction/Detail.html:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vaction/Detail.html:8
func StreamDetail(qw422016 *qt422016.Writer, cfg util.ValueMap, res *action.Result, includeSkipped bool, as *app.State, ps *cutil.PageState) {
//line views/vaction/Detail.html:8
	qw422016.N().S(`
  <div class="card">
    <h3>Config</h3>
    <table>
      <tbody>
`)
//line views/vaction/Detail.html:13
	for _, k := range cfg.Keys() {
//line views/vaction/Detail.html:13
		qw422016.N().S(`        <tr>
          <th class="shrink">`)
//line views/vaction/Detail.html:15
		qw422016.E().S(k)
//line views/vaction/Detail.html:15
		qw422016.N().S(`</th>
          <td>`)
//line views/vaction/Detail.html:16
		qw422016.E().V(cfg[k])
//line views/vaction/Detail.html:16
		qw422016.N().S(`</td>
        </tr>
`)
//line views/vaction/Detail.html:18
	}
//line views/vaction/Detail.html:18
	qw422016.N().S(`      </tbody>
    </table>
  </div>
`)
//line views/vaction/Detail.html:22
	if len(res.Errors) > 0 {
//line views/vaction/Detail.html:22
		qw422016.N().S(`  <div class="card">
    <div class="right">`)
//line views/vaction/Detail.html:24
		qw422016.N().D(len(res.Errors))
//line views/vaction/Detail.html:24
		qw422016.N().S(` `)
//line views/vaction/Detail.html:24
		qw422016.E().S(util.StringPluralMaybe("error", len(res.Errors)))
//line views/vaction/Detail.html:24
		qw422016.N().S(`</div>
    <h3>Errors</h3>
    <ul>
`)
//line views/vaction/Detail.html:27
		for _, e := range res.Errors {
//line views/vaction/Detail.html:27
			qw422016.N().S(`      <li class="error">`)
//line views/vaction/Detail.html:28
			qw422016.E().S(e)
//line views/vaction/Detail.html:28
			qw422016.N().S(`</li>
`)
//line views/vaction/Detail.html:29
		}
//line views/vaction/Detail.html:29
		qw422016.N().S(`    </ul>
  </div>
`)
//line views/vaction/Detail.html:32
	}
//line views/vaction/Detail.html:33
	if len(res.Logs) > 0 {
//line views/vaction/Detail.html:33
		qw422016.N().S(`  <div class="card">
    <h3>Logs</h3>
    <table>
      <tbody>
`)
//line views/vaction/Detail.html:38
		for _, l := range res.Logs {
//line views/vaction/Detail.html:38
			qw422016.N().S(`        <tr>
          <td><pre>`)
//line views/vaction/Detail.html:40
			qw422016.E().S(l)
//line views/vaction/Detail.html:40
			qw422016.N().S(`</pre></td>
        </tr>
`)
//line views/vaction/Detail.html:42
		}
//line views/vaction/Detail.html:42
		qw422016.N().S(`      </tbody>
    </table>
  </div>
`)
//line views/vaction/Detail.html:46
	}
//line views/vaction/Detail.html:47
	for _, mr := range res.Modules {
//line views/vaction/Detail.html:49
		links := ""
		for mIdx, mk := range mr.Keys {
			if mIdx > 0 {
				links += ", "
			}
			links += `<a href="/m/` + mk + `">` + mk + `</a>`
		}

//line views/vaction/Detail.html:56
		qw422016.N().S(`    <div class="card">
      <div class="right">`)
//line views/vaction/Detail.html:58
		qw422016.E().S(util.MicrosToMillis(mr.Duration))
//line views/vaction/Detail.html:58
		qw422016.N().S(`</div>
      <h3>`)
//line views/vaction/Detail.html:59
		qw422016.E().S(util.StringPluralMaybe("Module", len(mr.Keys)))
//line views/vaction/Detail.html:59
		qw422016.N().S(` [`)
//line views/vaction/Detail.html:59
		qw422016.N().S(links)
//line views/vaction/Detail.html:59
		qw422016.N().S(`]</h3>
      <div class="right">`)
//line views/vaction/Detail.html:60
		qw422016.N().S(res.StatusLog())
//line views/vaction/Detail.html:60
		qw422016.N().S(`</div>
      <em>`)
//line views/vaction/Detail.html:61
		qw422016.E().S(mr.Status)
//line views/vaction/Detail.html:61
		qw422016.N().S(`</em>
`)
//line views/vaction/Detail.html:62
		if len(mr.Actions) > 0 {
//line views/vaction/Detail.html:62
			qw422016.N().S(`      <h4>Actions</h4>
`)
//line views/vaction/Detail.html:64
			for _, a := range mr.Actions {
//line views/vaction/Detail.html:64
				qw422016.N().S(`        <a href="`)
//line views/vaction/Detail.html:65
				qw422016.E().S(a.URL())
//line views/vaction/Detail.html:65
				qw422016.N().S(`"><button>`)
//line views/vaction/Detail.html:65
				qw422016.E().S(a.Title)
//line views/vaction/Detail.html:65
				qw422016.N().S(`</button></a>
`)
//line views/vaction/Detail.html:66
			}
//line views/vaction/Detail.html:67
		}
//line views/vaction/Detail.html:68
		diffs := mr.DiffsFiltered(includeSkipped)

//line views/vaction/Detail.html:69
		if len(diffs) > 0 {
//line views/vaction/Detail.html:69
			qw422016.N().S(`      `)
//line views/vaction/Detail.html:70
			streamrenderDiffs(qw422016, res.Project.Key, diffs, as, ps)
//line views/vaction/Detail.html:70
			qw422016.N().S(`
`)
//line views/vaction/Detail.html:71
		}
//line views/vaction/Detail.html:71
		qw422016.N().S(`    </div>
`)
//line views/vaction/Detail.html:73
	}
//line views/vaction/Detail.html:74
}

//line views/vaction/Detail.html:74
func WriteDetail(qq422016 qtio422016.Writer, cfg util.ValueMap, res *action.Result, includeSkipped bool, as *app.State, ps *cutil.PageState) {
//line views/vaction/Detail.html:74
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vaction/Detail.html:74
	StreamDetail(qw422016, cfg, res, includeSkipped, as, ps)
//line views/vaction/Detail.html:74
	qt422016.ReleaseWriter(qw422016)
//line views/vaction/Detail.html:74
}

//line views/vaction/Detail.html:74
func Detail(cfg util.ValueMap, res *action.Result, includeSkipped bool, as *app.State, ps *cutil.PageState) string {
//line views/vaction/Detail.html:74
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vaction/Detail.html:74
	WriteDetail(qb422016, cfg, res, includeSkipped, as, ps)
//line views/vaction/Detail.html:74
	qs422016 := string(qb422016.B)
//line views/vaction/Detail.html:74
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vaction/Detail.html:74
	return qs422016
//line views/vaction/Detail.html:74
}