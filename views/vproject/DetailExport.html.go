// Code generated by qtc from "DetailExport.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vproject/DetailExport.html:1
package vproject

//line views/vproject/DetailExport.html:1
import (
	"fmt"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/components/view"
	"projectforge.dev/projectforge/views/vexport"
)

//line views/vproject/DetailExport.html:13
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vproject/DetailExport.html:13
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vproject/DetailExport.html:13
func StreamDetailExport(qw422016 *qt422016.Writer, key string, ea *metamodel.Args, as *app.State, ps *cutil.PageState) {
//line views/vproject/DetailExport.html:13
	qw422016.N().S(`
  <a href="/p/`)
//line views/vproject/DetailExport.html:14
	qw422016.E().S(key)
//line views/vproject/DetailExport.html:14
	qw422016.N().S(`/export"><button>`)
//line views/vproject/DetailExport.html:14
	components.StreamSVGButton(qw422016, "wrench", ps)
//line views/vproject/DetailExport.html:14
	qw422016.N().S(` Manage</button></a>
  <hr />
`)
//line views/vproject/DetailExport.html:16
	if ea.Config != nil && util.ToJSONCompact(ea.Config) != "{}" {
//line views/vproject/DetailExport.html:16
		qw422016.N().S(`  <h3>`)
//line views/vproject/DetailExport.html:17
		components.StreamSVGIcon(qw422016, `print`, ps)
//line views/vproject/DetailExport.html:17
		qw422016.N().S(` Export Configuration</h3>
  <div class="clear"></div>
  `)
//line views/vproject/DetailExport.html:19
		view.StreamMap(qw422016, true, ea.Config, ps)
//line views/vproject/DetailExport.html:19
		qw422016.N().S(`
  <hr />
`)
//line views/vproject/DetailExport.html:21
	}
//line views/vproject/DetailExport.html:22
	if len(ea.Groups) > 0 {
//line views/vproject/DetailExport.html:22
		qw422016.N().S(`  <h3>`)
//line views/vproject/DetailExport.html:23
		components.StreamSVGIcon(qw422016, `users`, ps)
//line views/vproject/DetailExport.html:23
		qw422016.N().S(` `)
//line views/vproject/DetailExport.html:23
		qw422016.E().S(util.StringPlural(len(ea.Groups), "Group"))
//line views/vproject/DetailExport.html:23
		qw422016.N().S(`</h3>
  <div class="mt clear"></div>
  `)
//line views/vproject/DetailExport.html:25
		vexport.StreamGroupList(qw422016, ea.Groups, 4)
//line views/vproject/DetailExport.html:25
		qw422016.N().S(`
  <hr />
`)
//line views/vproject/DetailExport.html:27
	}
//line views/vproject/DetailExport.html:28
	if len(ea.Enums) > 0 {
//line views/vproject/DetailExport.html:28
		qw422016.N().S(`  <h3>`)
//line views/vproject/DetailExport.html:29
		components.StreamSVGIcon(qw422016, `hammer`, ps)
//line views/vproject/DetailExport.html:29
		qw422016.N().S(` `)
//line views/vproject/DetailExport.html:29
		qw422016.E().S(util.StringPlural(len(ea.Enums), "Enum"))
//line views/vproject/DetailExport.html:29
		qw422016.N().S(`</h3>
  <div class="mt clear"></div>
  `)
//line views/vproject/DetailExport.html:31
		vexport.StreamEnumList(qw422016, ea.Enums, fmt.Sprintf("/p/%s/export/enums", key), as, ps)
//line views/vproject/DetailExport.html:31
		qw422016.N().S(`
  <hr />
`)
//line views/vproject/DetailExport.html:33
	}
//line views/vproject/DetailExport.html:34
	if len(ea.Models) > 0 {
//line views/vproject/DetailExport.html:34
		qw422016.N().S(`  <h3 class="mt">`)
//line views/vproject/DetailExport.html:35
		components.StreamSVGIcon(qw422016, `list`, ps)
//line views/vproject/DetailExport.html:35
		qw422016.N().S(` `)
//line views/vproject/DetailExport.html:35
		qw422016.E().S(util.StringPlural(len(ea.Models), "Model"))
//line views/vproject/DetailExport.html:35
		qw422016.N().S(`</h3>
  <div class="clear"></div>
  `)
//line views/vproject/DetailExport.html:37
		vexport.StreamModelList(qw422016, ea.Models, fmt.Sprintf("/p/%s/export/models", key), as, ps)
//line views/vproject/DetailExport.html:37
		qw422016.N().S(`
`)
//line views/vproject/DetailExport.html:38
	}
//line views/vproject/DetailExport.html:39
}

//line views/vproject/DetailExport.html:39
func WriteDetailExport(qq422016 qtio422016.Writer, key string, ea *metamodel.Args, as *app.State, ps *cutil.PageState) {
//line views/vproject/DetailExport.html:39
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vproject/DetailExport.html:39
	StreamDetailExport(qw422016, key, ea, as, ps)
//line views/vproject/DetailExport.html:39
	qt422016.ReleaseWriter(qw422016)
//line views/vproject/DetailExport.html:39
}

//line views/vproject/DetailExport.html:39
func DetailExport(key string, ea *metamodel.Args, as *app.State, ps *cutil.PageState) string {
//line views/vproject/DetailExport.html:39
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vproject/DetailExport.html:39
	WriteDetailExport(qb422016, key, ea, as, ps)
//line views/vproject/DetailExport.html:39
	qs422016 := string(qb422016.B)
//line views/vproject/DetailExport.html:39
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vproject/DetailExport.html:39
	return qs422016
//line views/vproject/DetailExport.html:39
}
