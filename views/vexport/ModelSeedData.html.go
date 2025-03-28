// Code generated by qtc from "ModelSeedData.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vexport/ModelSeedData.html:1
package vexport

//line views/vexport/ModelSeedData.html:1
import (
	"fmt"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vexport/ModelSeedData.html:11
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vexport/ModelSeedData.html:11
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vexport/ModelSeedData.html:11
type ModelSeedData struct {
	layout.Basic
	Model *model.Model
}

//line views/vexport/ModelSeedData.html:16
func (p *ModelSeedData) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vexport/ModelSeedData.html:16
	qw422016.N().S(`
`)
//line views/vexport/ModelSeedData.html:17
	m := p.Model

//line views/vexport/ModelSeedData.html:17
	qw422016.N().S(`  <div class="card">
    <h3>`)
//line views/vexport/ModelSeedData.html:19
	components.StreamSVGIcon(qw422016, m.IconSafe(), ps)
//line views/vexport/ModelSeedData.html:19
	qw422016.E().S(m.Name)
//line views/vexport/ModelSeedData.html:19
	qw422016.N().S(` Seed Data</h3>
    <div class="mt overflow">
      <div class="overflow full-width">
        <table>
          <thead>
            <tr>
`)
//line views/vexport/ModelSeedData.html:25
	for _, x := range p.Model.Columns {
//line views/vexport/ModelSeedData.html:25
		qw422016.N().S(`              <th>`)
//line views/vexport/ModelSeedData.html:26
		qw422016.E().S(x.Name)
//line views/vexport/ModelSeedData.html:26
		qw422016.N().S(`</th>
`)
//line views/vexport/ModelSeedData.html:27
	}
//line views/vexport/ModelSeedData.html:27
	qw422016.N().S(`            </tr>
          </thead>
          <tbody>
`)
//line views/vexport/ModelSeedData.html:31
	for _, row := range p.Model.SeedData {
//line views/vexport/ModelSeedData.html:31
		qw422016.N().S(`            <tr>
`)
//line views/vexport/ModelSeedData.html:33
		for _, cell := range row {
//line views/vexport/ModelSeedData.html:33
			qw422016.N().S(`              <td>`)
//line views/vexport/ModelSeedData.html:34
			qw422016.E().S(fmt.Sprint(cell))
//line views/vexport/ModelSeedData.html:34
			qw422016.N().S(`</td>
`)
//line views/vexport/ModelSeedData.html:35
		}
//line views/vexport/ModelSeedData.html:35
		qw422016.N().S(`            </tr>
`)
//line views/vexport/ModelSeedData.html:37
	}
//line views/vexport/ModelSeedData.html:37
	qw422016.N().S(`          </tbody>
        </table>
      </div>
    </div>
  </div>
`)
//line views/vexport/ModelSeedData.html:43
}

//line views/vexport/ModelSeedData.html:43
func (p *ModelSeedData) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vexport/ModelSeedData.html:43
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vexport/ModelSeedData.html:43
	p.StreamBody(qw422016, as, ps)
//line views/vexport/ModelSeedData.html:43
	qt422016.ReleaseWriter(qw422016)
//line views/vexport/ModelSeedData.html:43
}

//line views/vexport/ModelSeedData.html:43
func (p *ModelSeedData) Body(as *app.State, ps *cutil.PageState) string {
//line views/vexport/ModelSeedData.html:43
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vexport/ModelSeedData.html:43
	p.WriteBody(qb422016, as, ps)
//line views/vexport/ModelSeedData.html:43
	qs422016 := string(qb422016.B)
//line views/vexport/ModelSeedData.html:43
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vexport/ModelSeedData.html:43
	return qs422016
//line views/vexport/ModelSeedData.html:43
}
