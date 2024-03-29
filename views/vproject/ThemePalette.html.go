// Code generated by qtc from "ThemePalette.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vproject/ThemePalette.html:1
package vproject

//line views/vproject/ThemePalette.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vtheme"
)

//line views/vproject/ThemePalette.html:9
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vproject/ThemePalette.html:9
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vproject/ThemePalette.html:9
type ThemePalette struct {
	layout.Basic
	Project string
	Icon    string
	Palette string
	Themes  theme.Themes
	Title   string
}

//line views/vproject/ThemePalette.html:18
func (p *ThemePalette) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vproject/ThemePalette.html:18
	qw422016.N().S(`
  <div class="card">
    <h3>Set theme for project [`)
//line views/vproject/ThemePalette.html:20
	qw422016.E().S(p.Project)
//line views/vproject/ThemePalette.html:20
	qw422016.N().S(`]</h3>
    <form action="/theme" method="post">
      <input type="hidden" name="palette" value="`)
//line views/vproject/ThemePalette.html:22
	qw422016.E().S(p.Palette)
//line views/vproject/ThemePalette.html:22
	qw422016.N().S(`" />
      <div class="overflow full-width">
        <table class="mt">
          <tbody>
`)
//line views/vproject/ThemePalette.html:26
	for _, t := range p.Themes {
//line views/vproject/ThemePalette.html:26
		qw422016.N().S(`            <tr>
`)
//line views/vproject/ThemePalette.html:28
		if p.Project == "" {
//line views/vproject/ThemePalette.html:28
			qw422016.N().S(`              <th class="shrink"><input type="checkbox" id="`)
//line views/vproject/ThemePalette.html:29
			qw422016.E().S(t.Key)
//line views/vproject/ThemePalette.html:29
			qw422016.N().S(`" name="`)
//line views/vproject/ThemePalette.html:29
			qw422016.E().S(t.Key)
//line views/vproject/ThemePalette.html:29
			qw422016.N().S(`" value="true" /></th>
`)
//line views/vproject/ThemePalette.html:30
		}
//line views/vproject/ThemePalette.html:30
		qw422016.N().S(`              <th>
                <label for="`)
//line views/vproject/ThemePalette.html:32
		qw422016.E().S(t.Key)
//line views/vproject/ThemePalette.html:32
		qw422016.N().S(`">`)
//line views/vproject/ThemePalette.html:32
		qw422016.E().S(t.Key)
//line views/vproject/ThemePalette.html:32
		qw422016.N().S(`</label>
                <a href="/p/`)
//line views/vproject/ThemePalette.html:33
		qw422016.E().S(p.Project)
//line views/vproject/ThemePalette.html:33
		qw422016.N().S(`/palette/`)
//line views/vproject/ThemePalette.html:33
		qw422016.E().S(p.Palette)
//line views/vproject/ThemePalette.html:33
		qw422016.N().S(`/`)
//line views/vproject/ThemePalette.html:33
		qw422016.E().S(t.Key)
//line views/vproject/ThemePalette.html:33
		qw422016.N().S(`" class="link-confirm" data-message="Are you sure you'd like to overwrite the default theme?">Set as Default</a>
              </th>
              <th class="shrink" style="background-color: #ffffff; padding: 12px 36px;">`)
//line views/vproject/ThemePalette.html:35
		vtheme.StreamMockupColors(qw422016, p.Title, "", t.Light, true, p.Icon, 5, ps)
//line views/vproject/ThemePalette.html:35
		qw422016.N().S(`</th>
              <th class="shrink" style="background-color: #121212; padding: 12px 36px;">`)
//line views/vproject/ThemePalette.html:36
		vtheme.StreamMockupColors(qw422016, p.Title, "", t.Dark, true, p.Icon, 5, ps)
//line views/vproject/ThemePalette.html:36
		qw422016.N().S(`</th>
            </tr>
`)
//line views/vproject/ThemePalette.html:38
	}
//line views/vproject/ThemePalette.html:38
	qw422016.N().S(`          </tbody>
        </table>
      </div>
    </form>
  </div>
`)
//line views/vproject/ThemePalette.html:44
}

//line views/vproject/ThemePalette.html:44
func (p *ThemePalette) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vproject/ThemePalette.html:44
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vproject/ThemePalette.html:44
	p.StreamBody(qw422016, as, ps)
//line views/vproject/ThemePalette.html:44
	qt422016.ReleaseWriter(qw422016)
//line views/vproject/ThemePalette.html:44
}

//line views/vproject/ThemePalette.html:44
func (p *ThemePalette) Body(as *app.State, ps *cutil.PageState) string {
//line views/vproject/ThemePalette.html:44
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vproject/ThemePalette.html:44
	p.WriteBody(qb422016, as, ps)
//line views/vproject/ThemePalette.html:44
	qs422016 := string(qb422016.B)
//line views/vproject/ThemePalette.html:44
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vproject/ThemePalette.html:44
	return qs422016
//line views/vproject/ThemePalette.html:44
}
