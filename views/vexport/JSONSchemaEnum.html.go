// Code generated by qtc from "JSONSchemaEnum.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vexport/JSONSchemaEnum.html:1
package vexport

//line views/vexport/JSONSchemaEnum.html:1
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vexport/JSONSchemaEnum.html:12
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vexport/JSONSchemaEnum.html:12
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vexport/JSONSchemaEnum.html:12
type JSONSchemaEnum struct {
	layout.Basic
	Project    *project.Project
	Enum       *enum.Enum
	Collection *jsonschema.Collection
	Result     *enum.Enum
}

//line views/vexport/JSONSchemaEnum.html:20
func (p *JSONSchemaEnum) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vexport/JSONSchemaEnum.html:20
	qw422016.N().S(`
`)
//line views/vexport/JSONSchemaEnum.html:21
	x := p.Enum

//line views/vexport/JSONSchemaEnum.html:21
	qw422016.N().S(`  <div class="card">
`)
//line views/vexport/JSONSchemaEnum.html:23
	sch := p.Collection.GetSchema(x.ID())

//line views/vexport/JSONSchemaEnum.html:24
	df := util.DiffObjects(x, p.Result)

//line views/vexport/JSONSchemaEnum.html:25
	if len(df) > 0 {
//line views/vexport/JSONSchemaEnum.html:25
		qw422016.N().S(`    <div class="right">
      `)
//line views/vexport/JSONSchemaEnum.html:27
		qw422016.E().S(util.StringPlural(len(df), "difference"))
//line views/vexport/JSONSchemaEnum.html:27
		qw422016.N().S(`
    </div>
`)
//line views/vexport/JSONSchemaEnum.html:29
	}
//line views/vexport/JSONSchemaEnum.html:29
	qw422016.N().S(`    <h3>`)
//line views/vexport/JSONSchemaEnum.html:30
	components.StreamSVGIcon(qw422016, `table`, ps)
//line views/vexport/JSONSchemaEnum.html:30
	qw422016.N().S(` [`)
//line views/vexport/JSONSchemaEnum.html:30
	qw422016.E().S(p.Project.Key)
//line views/vexport/JSONSchemaEnum.html:30
	qw422016.N().S(` / `)
//line views/vexport/JSONSchemaEnum.html:30
	qw422016.E().S(x.Name)
//line views/vexport/JSONSchemaEnum.html:30
	qw422016.N().S(`] JSON Schema</h3>
    `)
//line views/vexport/JSONSchemaEnum.html:31
	streamrenderJSONSchemaEnum(qw422016, p.Project, x, sch, p.Result, df, ps)
//line views/vexport/JSONSchemaEnum.html:31
	qw422016.N().S(`
  </div>
`)
//line views/vexport/JSONSchemaEnum.html:33
}

//line views/vexport/JSONSchemaEnum.html:33
func (p *JSONSchemaEnum) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vexport/JSONSchemaEnum.html:33
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vexport/JSONSchemaEnum.html:33
	p.StreamBody(qw422016, as, ps)
//line views/vexport/JSONSchemaEnum.html:33
	qt422016.ReleaseWriter(qw422016)
//line views/vexport/JSONSchemaEnum.html:33
}

//line views/vexport/JSONSchemaEnum.html:33
func (p *JSONSchemaEnum) Body(as *app.State, ps *cutil.PageState) string {
//line views/vexport/JSONSchemaEnum.html:33
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vexport/JSONSchemaEnum.html:33
	p.WriteBody(qb422016, as, ps)
//line views/vexport/JSONSchemaEnum.html:33
	qs422016 := string(qb422016.B)
//line views/vexport/JSONSchemaEnum.html:33
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vexport/JSONSchemaEnum.html:33
	return qs422016
//line views/vexport/JSONSchemaEnum.html:33
}

//line views/vexport/JSONSchemaEnum.html:35
func streamrenderJSONSchemaEnum(qw422016 *qt422016.Writer, prj *project.Project, x *enum.Enum, sch *jsonschema.Schema, result *enum.Enum, df util.Diffs, ps *cutil.PageState) {
//line views/vexport/JSONSchemaEnum.html:35
	qw422016.N().S(`
  <div class="flex" style="">
    <div class="flex-item">
      <a href="`)
//line views/vexport/JSONSchemaEnum.html:38
	qw422016.E().S(prj.WebPath())
//line views/vexport/JSONSchemaEnum.html:38
	qw422016.N().S(`/export/enums/`)
//line views/vexport/JSONSchemaEnum.html:38
	qw422016.E().S(x.Name)
//line views/vexport/JSONSchemaEnum.html:38
	qw422016.N().S(`"><strong>Original</strong></a>
      <pre class="mt">`)
//line views/vexport/JSONSchemaEnum.html:39
	qw422016.E().S(util.ToJSON(x))
//line views/vexport/JSONSchemaEnum.html:39
	qw422016.N().S(`</pre>
    </div>
    <div class="flex-item">
      <strong>Result</strong>
      <pre class="mt">`)
//line views/vexport/JSONSchemaEnum.html:43
	qw422016.E().S(util.ToJSON(result))
//line views/vexport/JSONSchemaEnum.html:43
	qw422016.N().S(`</pre>
    </div>
  </div>
  <div class="flex" style="">
    <div class="flex-item">
      <strong>Schema</strong>
      <pre class="mt">`)
//line views/vexport/JSONSchemaEnum.html:49
	qw422016.E().S(util.ToJSON(sch))
//line views/vexport/JSONSchemaEnum.html:49
	qw422016.N().S(`</pre>
    </div>
    <div class="flex-item">
      <strong>Diff</strong>
      <pre class="mt">`)
//line views/vexport/JSONSchemaEnum.html:53
	qw422016.E().S(util.ToJSON(df))
//line views/vexport/JSONSchemaEnum.html:53
	qw422016.N().S(`</pre>
    </div>
  </div>
`)
//line views/vexport/JSONSchemaEnum.html:56
}

//line views/vexport/JSONSchemaEnum.html:56
func writerenderJSONSchemaEnum(qq422016 qtio422016.Writer, prj *project.Project, x *enum.Enum, sch *jsonschema.Schema, result *enum.Enum, df util.Diffs, ps *cutil.PageState) {
//line views/vexport/JSONSchemaEnum.html:56
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vexport/JSONSchemaEnum.html:56
	streamrenderJSONSchemaEnum(qw422016, prj, x, sch, result, df, ps)
//line views/vexport/JSONSchemaEnum.html:56
	qt422016.ReleaseWriter(qw422016)
//line views/vexport/JSONSchemaEnum.html:56
}

//line views/vexport/JSONSchemaEnum.html:56
func renderJSONSchemaEnum(prj *project.Project, x *enum.Enum, sch *jsonschema.Schema, result *enum.Enum, df util.Diffs, ps *cutil.PageState) string {
//line views/vexport/JSONSchemaEnum.html:56
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vexport/JSONSchemaEnum.html:56
	writerenderJSONSchemaEnum(qb422016, prj, x, sch, result, df, ps)
//line views/vexport/JSONSchemaEnum.html:56
	qs422016 := string(qb422016.B)
//line views/vexport/JSONSchemaEnum.html:56
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vexport/JSONSchemaEnum.html:56
	return qs422016
//line views/vexport/JSONSchemaEnum.html:56
}
