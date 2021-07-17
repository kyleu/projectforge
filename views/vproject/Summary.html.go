// Code generated by qtc from "Summary.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vproject/Summary.html:1
package vproject

//line views/vproject/Summary.html:1
import (
	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/views/components"
)

//line views/vproject/Summary.html:7
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vproject/Summary.html:7
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vproject/Summary.html:7
func StreamSummary(qw422016 *qt422016.Writer, prj *project.Project) {
//line views/vproject/Summary.html:7
	qw422016.N().S(`
  <div class="card">
    <div class="right"><a href="#modal-project"><button type="button">JSON</button></a></div>
    <h3>`)
//line views/vproject/Summary.html:10
	qw422016.E().S(prj.Name)
//line views/vproject/Summary.html:10
	qw422016.N().S(`</h3>
    <div class="mt">
`)
//line views/vproject/Summary.html:12
	for _, t := range action.ProjectTypes {
//line views/vproject/Summary.html:12
		qw422016.N().S(`      <a href="/run/`)
//line views/vproject/Summary.html:13
		qw422016.E().S(prj.Key)
//line views/vproject/Summary.html:13
		qw422016.N().S(`/`)
//line views/vproject/Summary.html:13
		qw422016.E().S(t.Key)
//line views/vproject/Summary.html:13
		qw422016.N().S(`" title="`)
//line views/vproject/Summary.html:13
		qw422016.E().S(t.Description)
//line views/vproject/Summary.html:13
		qw422016.N().S(`"><button>`)
//line views/vproject/Summary.html:13
		qw422016.E().S(t.Title)
//line views/vproject/Summary.html:13
		qw422016.N().S(`</button></a>
`)
//line views/vproject/Summary.html:14
	}
//line views/vproject/Summary.html:14
	qw422016.N().S(`    </div>
  </div>
  `)
//line views/vproject/Summary.html:17
	components.StreamJSONModal(qw422016, "project", "Project JSON", prj, 1)
//line views/vproject/Summary.html:17
	qw422016.N().S(`
`)
//line views/vproject/Summary.html:18
}

//line views/vproject/Summary.html:18
func WriteSummary(qq422016 qtio422016.Writer, prj *project.Project) {
//line views/vproject/Summary.html:18
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vproject/Summary.html:18
	StreamSummary(qw422016, prj)
//line views/vproject/Summary.html:18
	qt422016.ReleaseWriter(qw422016)
//line views/vproject/Summary.html:18
}

//line views/vproject/Summary.html:18
func Summary(prj *project.Project) string {
//line views/vproject/Summary.html:18
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vproject/Summary.html:18
	WriteSummary(qb422016, prj)
//line views/vproject/Summary.html:18
	qs422016 := string(qb422016.B)
//line views/vproject/Summary.html:18
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vproject/Summary.html:18
	return qs422016
//line views/vproject/Summary.html:18
}