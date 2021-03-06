// Code generated by qtc from "Add.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vtheme/Add.html:1
package vtheme

//line views/vtheme/Add.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vtheme/Add.html:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vtheme/Add.html:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vtheme/Add.html:8
type Add struct {
	layout.Basic
	Project string
	Palette string
	Themes  theme.Themes
	Title   string
}

//line views/vtheme/Add.html:16
func (p *Add) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vtheme/Add.html:16
	qw422016.N().S(`
  <div class="card">
`)
//line views/vtheme/Add.html:18
	if p.Project == "" {
//line views/vtheme/Add.html:18
		qw422016.N().S(`    <h3>Add Themes</h3>
`)
//line views/vtheme/Add.html:20
	} else {
//line views/vtheme/Add.html:20
		qw422016.N().S(`    <h3>Set theme for project [`)
//line views/vtheme/Add.html:21
		qw422016.E().S(p.Project)
//line views/vtheme/Add.html:21
		qw422016.N().S(`]</h3>
`)
//line views/vtheme/Add.html:22
	}
//line views/vtheme/Add.html:22
	qw422016.N().S(`    <form action="/theme" method="post">
      <input type="hidden" name="palette" value="`)
//line views/vtheme/Add.html:24
	qw422016.E().S(p.Palette)
//line views/vtheme/Add.html:24
	qw422016.N().S(`" />
      <table class="mt">
        <tbody>
`)
//line views/vtheme/Add.html:27
	for _, t := range p.Themes {
//line views/vtheme/Add.html:27
		qw422016.N().S(`          <tr>
`)
//line views/vtheme/Add.html:29
		if p.Project == "" {
//line views/vtheme/Add.html:29
			qw422016.N().S(`            <th class="shrink"><input type="checkbox" id="`)
//line views/vtheme/Add.html:30
			qw422016.E().S(t.Key)
//line views/vtheme/Add.html:30
			qw422016.N().S(`" name="`)
//line views/vtheme/Add.html:30
			qw422016.E().S(t.Key)
//line views/vtheme/Add.html:30
			qw422016.N().S(`" value="true" /></th>
`)
//line views/vtheme/Add.html:31
		}
//line views/vtheme/Add.html:31
		qw422016.N().S(`            <th>
              <label for="`)
//line views/vtheme/Add.html:33
		qw422016.E().S(t.Key)
//line views/vtheme/Add.html:33
		qw422016.N().S(`">`)
//line views/vtheme/Add.html:33
		qw422016.E().S(t.Key)
//line views/vtheme/Add.html:33
		qw422016.N().S(`</label>
`)
//line views/vtheme/Add.html:34
		if p.Project == "" {
//line views/vtheme/Add.html:34
			qw422016.N().S(`              <a href="/theme/preview/`)
//line views/vtheme/Add.html:35
			qw422016.E().S(p.Palette)
//line views/vtheme/Add.html:35
			qw422016.N().S(`/`)
//line views/vtheme/Add.html:35
			qw422016.E().S(t.Key)
//line views/vtheme/Add.html:35
			qw422016.N().S(`">Preview</a>
`)
//line views/vtheme/Add.html:36
		} else {
//line views/vtheme/Add.html:36
			qw422016.N().S(`              <a href="/theme/preview/`)
//line views/vtheme/Add.html:37
			qw422016.E().S(p.Palette)
//line views/vtheme/Add.html:37
			qw422016.N().S(`/`)
//line views/vtheme/Add.html:37
			qw422016.E().S(t.Key)
//line views/vtheme/Add.html:37
			qw422016.N().S(`?project=`)
//line views/vtheme/Add.html:37
			qw422016.E().S(p.Project)
//line views/vtheme/Add.html:37
			qw422016.N().S(`">Set</a>
`)
//line views/vtheme/Add.html:38
		}
//line views/vtheme/Add.html:38
		qw422016.N().S(`            </th>
            <th class="shrink" style="background-color: #ffffff; padding: 12px 36px;">`)
//line views/vtheme/Add.html:40
		StreamMockupColors(qw422016, p.Title, "", t.Light, true, 5, ps)
//line views/vtheme/Add.html:40
		qw422016.N().S(`</th>
            <th class="shrink" style="background-color: #121212; padding: 12px 36px;">`)
//line views/vtheme/Add.html:41
		StreamMockupColors(qw422016, p.Title, "", t.Dark, true, 5, ps)
//line views/vtheme/Add.html:41
		qw422016.N().S(`</th>
          </tr>
`)
//line views/vtheme/Add.html:43
	}
//line views/vtheme/Add.html:43
	qw422016.N().S(`        </tbody>
      </table>
    </form>
  </div>
`)
//line views/vtheme/Add.html:48
}

//line views/vtheme/Add.html:48
func (p *Add) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vtheme/Add.html:48
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vtheme/Add.html:48
	p.StreamBody(qw422016, as, ps)
//line views/vtheme/Add.html:48
	qt422016.ReleaseWriter(qw422016)
//line views/vtheme/Add.html:48
}

//line views/vtheme/Add.html:48
func (p *Add) Body(as *app.State, ps *cutil.PageState) string {
//line views/vtheme/Add.html:48
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vtheme/Add.html:48
	p.WriteBody(qb422016, as, ps)
//line views/vtheme/Add.html:48
	qs422016 := string(qb422016.B)
//line views/vtheme/Add.html:48
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vtheme/Add.html:48
	return qs422016
//line views/vtheme/Add.html:48
}
