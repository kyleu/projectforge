// Code generated by qtc from "Overview.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vexport/Overview.html:1
package vexport

//line views/vexport/Overview.html:1
import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vexport/Overview.html:15
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vexport/Overview.html:15
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vexport/Overview.html:15
type Overview struct {
	layout.Basic
	Project *project.Project
	Args    *model.Args
}

//line views/vexport/Overview.html:21
func (p *Overview) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vexport/Overview.html:21
	qw422016.N().S(`
  <div class="card">
    <div class="right"><a href="/p/`)
//line views/vexport/Overview.html:23
	qw422016.E().S(p.Project.Key)
//line views/vexport/Overview.html:23
	qw422016.N().S(`/export/config"><button>`)
//line views/vexport/Overview.html:23
	components.StreamSVGButton(qw422016, "edit", ps)
//line views/vexport/Overview.html:23
	qw422016.N().S(` Edit</button></a></div>
    <h3>`)
//line views/vexport/Overview.html:24
	components.StreamSVGIcon(qw422016, `print`, ps)
//line views/vexport/Overview.html:24
	qw422016.N().S(` Export Configuration</h3>
    `)
//line views/vexport/Overview.html:25
	components.StreamJSON(qw422016, p.Args.Config)
//line views/vexport/Overview.html:25
	qw422016.N().S(`
  </div>
  <div class="card">
    <div class="right"><a href="/p/`)
//line views/vexport/Overview.html:28
	qw422016.E().S(p.Project.Key)
//line views/vexport/Overview.html:28
	qw422016.N().S(`/export/groups"><button>`)
//line views/vexport/Overview.html:28
	components.StreamSVGButton(qw422016, "edit", ps)
//line views/vexport/Overview.html:28
	qw422016.N().S(` Edit</button></a></div>
    <h3>`)
//line views/vexport/Overview.html:29
	components.StreamSVGIcon(qw422016, `users`, ps)
//line views/vexport/Overview.html:29
	qw422016.N().S(` `)
//line views/vexport/Overview.html:29
	qw422016.E().S(util.StringPlural(len(p.Args.Groups), "Group"))
//line views/vexport/Overview.html:29
	qw422016.N().S(`</h3>
    <div class="mt">
`)
//line views/vexport/Overview.html:31
	if len(p.Args.Groups) == 0 {
//line views/vexport/Overview.html:31
		qw422016.N().S(`      <em>no groups defined</em>
`)
//line views/vexport/Overview.html:33
	} else {
//line views/vexport/Overview.html:33
		qw422016.N().S(`      `)
//line views/vexport/Overview.html:34
		StreamGroupList(qw422016, p.Args.Groups, 2)
//line views/vexport/Overview.html:34
		qw422016.N().S(`
`)
//line views/vexport/Overview.html:35
	}
//line views/vexport/Overview.html:35
	qw422016.N().S(`    </div>
  </div>
  <div class="card">
    <div class="right">
      <a href="/p/`)
//line views/vexport/Overview.html:40
	qw422016.E().S(p.Project.Key)
//line views/vexport/Overview.html:40
	qw422016.N().S(`/export/enums/create/new"><button>`)
//line views/vexport/Overview.html:40
	components.StreamSVGButton(qw422016, "plus", ps)
//line views/vexport/Overview.html:40
	qw422016.N().S(` New</button></a>
    </div>
    <h3>`)
//line views/vexport/Overview.html:42
	components.StreamSVGIcon(qw422016, `hammer`, ps)
//line views/vexport/Overview.html:42
	qw422016.N().S(` `)
//line views/vexport/Overview.html:42
	qw422016.E().S(util.StringPlural(len(p.Args.Enums), "Enum"))
//line views/vexport/Overview.html:42
	qw422016.N().S(`</h3>
`)
//line views/vexport/Overview.html:43
	if len(p.Args.Enums) == 0 {
//line views/vexport/Overview.html:43
		qw422016.N().S(`    <em>no enums defined</em>
`)
//line views/vexport/Overview.html:45
	} else {
//line views/vexport/Overview.html:45
		qw422016.N().S(`    `)
//line views/vexport/Overview.html:46
		StreamEnumList(qw422016, p.Args.Enums, fmt.Sprintf("/p/%s/export/enums", p.Project.Key), as, ps)
//line views/vexport/Overview.html:46
		qw422016.N().S(`
`)
//line views/vexport/Overview.html:47
	}
//line views/vexport/Overview.html:47
	qw422016.N().S(`  </div>
  <div class="card">
    <div class="right">
      <a href="/p/`)
//line views/vexport/Overview.html:51
	qw422016.E().S(p.Project.Key)
//line views/vexport/Overview.html:51
	qw422016.N().S(`/export/models/create/derive"><button>`)
//line views/vexport/Overview.html:51
	components.StreamSVGButton(qw422016, "dna", ps)
//line views/vexport/Overview.html:51
	qw422016.N().S(` Derive</button></a>
      <a href="/p/`)
//line views/vexport/Overview.html:52
	qw422016.E().S(p.Project.Key)
//line views/vexport/Overview.html:52
	qw422016.N().S(`/export/models/create/new"><button>`)
//line views/vexport/Overview.html:52
	components.StreamSVGButton(qw422016, "plus", ps)
//line views/vexport/Overview.html:52
	qw422016.N().S(` New</button></a>
    </div>
    <h3>`)
//line views/vexport/Overview.html:54
	components.StreamSVGIcon(qw422016, `list`, ps)
//line views/vexport/Overview.html:54
	qw422016.N().S(` `)
//line views/vexport/Overview.html:54
	qw422016.E().S(util.StringPlural(len(p.Args.Models), "Model"))
//line views/vexport/Overview.html:54
	qw422016.N().S(`</h3>
    `)
//line views/vexport/Overview.html:55
	StreamModelList(qw422016, p.Args.Models, fmt.Sprintf("/p/%s/export/models", p.Project.Key), as, ps)
//line views/vexport/Overview.html:55
	qw422016.N().S(`
  </div>
`)
//line views/vexport/Overview.html:57
}

//line views/vexport/Overview.html:57
func (p *Overview) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vexport/Overview.html:57
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vexport/Overview.html:57
	p.StreamBody(qw422016, as, ps)
//line views/vexport/Overview.html:57
	qt422016.ReleaseWriter(qw422016)
//line views/vexport/Overview.html:57
}

//line views/vexport/Overview.html:57
func (p *Overview) Body(as *app.State, ps *cutil.PageState) string {
//line views/vexport/Overview.html:57
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vexport/Overview.html:57
	p.WriteBody(qb422016, as, ps)
//line views/vexport/Overview.html:57
	qs422016 := string(qb422016.B)
//line views/vexport/Overview.html:57
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vexport/Overview.html:57
	return qs422016
//line views/vexport/Overview.html:57
}

//line views/vexport/Overview.html:59
func StreamGroupList(qw422016 *qt422016.Writer, groups model.Groups, indent int) {
//line views/vexport/Overview.html:60
	components.StreamIndent(qw422016, true, indent)
//line views/vexport/Overview.html:60
	qw422016.N().S(`<ul>`)
//line views/vexport/Overview.html:62
	for _, g := range groups {
//line views/vexport/Overview.html:63
		components.StreamIndent(qw422016, true, indent+1)
//line views/vexport/Overview.html:63
		qw422016.N().S(`<li>`)
//line views/vexport/Overview.html:65
		components.StreamIndent(qw422016, true, indent+2)
//line views/vexport/Overview.html:66
		qw422016.E().S(g.String())
//line views/vexport/Overview.html:67
		if g.Description != "" {
//line views/vexport/Overview.html:67
			qw422016.N().S(`:`)
//line views/vexport/Overview.html:68
			qw422016.N().S(` `)
//line views/vexport/Overview.html:68
			qw422016.E().S(g.Description)
//line views/vexport/Overview.html:69
		}
//line views/vexport/Overview.html:70
		if g.Route != "" {
//line views/vexport/Overview.html:70
			qw422016.N().S(`:`)
//line views/vexport/Overview.html:71
			qw422016.N().S(` `)
//line views/vexport/Overview.html:71
			qw422016.N().S(`<em><code>`)
//line views/vexport/Overview.html:71
			qw422016.E().S(g.Route)
//line views/vexport/Overview.html:71
			qw422016.N().S(`</code></em>`)
//line views/vexport/Overview.html:72
		}
//line views/vexport/Overview.html:73
		if len(g.Children) > 0 {
//line views/vexport/Overview.html:74
			StreamGroupList(qw422016, g.Children, indent+3)
//line views/vexport/Overview.html:75
		}
//line views/vexport/Overview.html:76
		components.StreamIndent(qw422016, true, indent+1)
//line views/vexport/Overview.html:76
		qw422016.N().S(`</li>`)
//line views/vexport/Overview.html:78
	}
//line views/vexport/Overview.html:79
	components.StreamIndent(qw422016, true, indent)
//line views/vexport/Overview.html:79
	qw422016.N().S(`</ul>`)
//line views/vexport/Overview.html:81
}

//line views/vexport/Overview.html:81
func WriteGroupList(qq422016 qtio422016.Writer, groups model.Groups, indent int) {
//line views/vexport/Overview.html:81
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vexport/Overview.html:81
	StreamGroupList(qw422016, groups, indent)
//line views/vexport/Overview.html:81
	qt422016.ReleaseWriter(qw422016)
//line views/vexport/Overview.html:81
}

//line views/vexport/Overview.html:81
func GroupList(groups model.Groups, indent int) string {
//line views/vexport/Overview.html:81
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vexport/Overview.html:81
	WriteGroupList(qb422016, groups, indent)
//line views/vexport/Overview.html:81
	qs422016 := string(qb422016.B)
//line views/vexport/Overview.html:81
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vexport/Overview.html:81
	return qs422016
//line views/vexport/Overview.html:81
}

//line views/vexport/Overview.html:83
func StreamEnumList(qw422016 *qt422016.Writer, enums enum.Enums, urlPrefix string, as *app.State, ps *cutil.PageState) {
//line views/vexport/Overview.html:83
	qw422016.N().S(`
  <div class="overflow full-width">
    <table class="mt min-200 full-width">
      <tbody>
`)
//line views/vexport/Overview.html:87
	for _, e := range enums {
//line views/vexport/Overview.html:89
		u := fmt.Sprintf("%s/%s", urlPrefix, e.Name)
		var prefix string
		if len(e.Group) > 0 {
			prefix += strings.Join(e.Group, "/") + ", "
		}

//line views/vexport/Overview.html:94
		qw422016.N().S(`        <tr>
          <td class="shrink"><a href="`)
//line views/vexport/Overview.html:96
		qw422016.E().S(u)
//line views/vexport/Overview.html:96
		qw422016.N().S(`">`)
//line views/vexport/Overview.html:96
		components.StreamSVGIcon(qw422016, e.IconSafe(), ps)
//line views/vexport/Overview.html:96
		qw422016.N().S(`</a> <a href="`)
//line views/vexport/Overview.html:96
		qw422016.E().S(u)
//line views/vexport/Overview.html:96
		qw422016.N().S(`">`)
//line views/vexport/Overview.html:96
		qw422016.E().S(e.Title())
//line views/vexport/Overview.html:96
		qw422016.N().S(`</a></td>
          <td>`)
//line views/vexport/Overview.html:97
		qw422016.E().S(e.Description)
//line views/vexport/Overview.html:97
		qw422016.N().S(`</td>
        </tr>
`)
//line views/vexport/Overview.html:99
	}
//line views/vexport/Overview.html:99
	qw422016.N().S(`      </tbody>
    </table>
  </div>
`)
//line views/vexport/Overview.html:103
}

//line views/vexport/Overview.html:103
func WriteEnumList(qq422016 qtio422016.Writer, enums enum.Enums, urlPrefix string, as *app.State, ps *cutil.PageState) {
//line views/vexport/Overview.html:103
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vexport/Overview.html:103
	StreamEnumList(qw422016, enums, urlPrefix, as, ps)
//line views/vexport/Overview.html:103
	qt422016.ReleaseWriter(qw422016)
//line views/vexport/Overview.html:103
}

//line views/vexport/Overview.html:103
func EnumList(enums enum.Enums, urlPrefix string, as *app.State, ps *cutil.PageState) string {
//line views/vexport/Overview.html:103
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vexport/Overview.html:103
	WriteEnumList(qb422016, enums, urlPrefix, as, ps)
//line views/vexport/Overview.html:103
	qs422016 := string(qb422016.B)
//line views/vexport/Overview.html:103
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vexport/Overview.html:103
	return qs422016
//line views/vexport/Overview.html:103
}

//line views/vexport/Overview.html:105
func StreamModelList(qw422016 *qt422016.Writer, models model.Models, urlPrefix string, as *app.State, ps *cutil.PageState) {
//line views/vexport/Overview.html:105
	qw422016.N().S(`
  <div class="overflow full-width">
    <table class="mt min-200 full-width">
      <tbody>
`)
//line views/vexport/Overview.html:109
	for _, m := range models {
//line views/vexport/Overview.html:111
		u := fmt.Sprintf("%s/%s", urlPrefix, m.Name)
		var prefix string
		if len(m.Group) > 0 {
			prefix += strings.Join(m.Group, "/") + ", "
		}
		if len(m.SeedData) > 0 {
			prefix += fmt.Sprintf("%s of seed data", util.StringPlural(len(m.SeedData), "row"))
		}

//line views/vexport/Overview.html:119
		qw422016.N().S(`        <tr>
          <td class="shrink"><a href="`)
//line views/vexport/Overview.html:121
		qw422016.E().S(u)
//line views/vexport/Overview.html:121
		qw422016.N().S(`">`)
//line views/vexport/Overview.html:121
		components.StreamSVGIcon(qw422016, m.IconSafe(), ps)
//line views/vexport/Overview.html:121
		qw422016.N().S(`</a> <a href="`)
//line views/vexport/Overview.html:121
		qw422016.E().S(u)
//line views/vexport/Overview.html:121
		qw422016.N().S(`">`)
//line views/vexport/Overview.html:121
		qw422016.E().S(m.Title())
//line views/vexport/Overview.html:121
		qw422016.N().S(`</a></td>
          <td class="text-align-right"><em>`)
//line views/vexport/Overview.html:122
		qw422016.E().S(prefix)
//line views/vexport/Overview.html:122
		qw422016.E().S(util.StringPlural(len(m.Columns), "field"))
//line views/vexport/Overview.html:122
		qw422016.N().S(`</em></td>
        </tr>
`)
//line views/vexport/Overview.html:124
	}
//line views/vexport/Overview.html:124
	qw422016.N().S(`      </tbody>
    </table>
  </div>
`)
//line views/vexport/Overview.html:128
}

//line views/vexport/Overview.html:128
func WriteModelList(qq422016 qtio422016.Writer, models model.Models, urlPrefix string, as *app.State, ps *cutil.PageState) {
//line views/vexport/Overview.html:128
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vexport/Overview.html:128
	StreamModelList(qw422016, models, urlPrefix, as, ps)
//line views/vexport/Overview.html:128
	qt422016.ReleaseWriter(qw422016)
//line views/vexport/Overview.html:128
}

//line views/vexport/Overview.html:128
func ModelList(models model.Models, urlPrefix string, as *app.State, ps *cutil.PageState) string {
//line views/vexport/Overview.html:128
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vexport/Overview.html:128
	WriteModelList(qb422016, models, urlPrefix, as, ps)
//line views/vexport/Overview.html:128
	qs422016 := string(qb422016.B)
//line views/vexport/Overview.html:128
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vexport/Overview.html:128
	return qs422016
//line views/vexport/Overview.html:128
}
