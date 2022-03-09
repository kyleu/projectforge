// Code generated by qtc from "Welcome.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vwelcome/Welcome.html:1
package vwelcome

//line views/vwelcome/Welcome.html:1
import (
	"strings"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vwelcome/Welcome.html:12
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vwelcome/Welcome.html:12
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vwelcome/Welcome.html:12
type Welcome struct {
	layout.Basic
	Project *project.Project
}

//line views/vwelcome/Welcome.html:17
func (p *Welcome) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vwelcome/Welcome.html:17
	qw422016.N().S(`
`)
//line views/vwelcome/Welcome.html:19
	prj := p.Project
	info := prj.Info
	if info == nil {
		info = &project.Info{}
	}

//line views/vwelcome/Welcome.html:24
	qw422016.N().S(`  <div class="card">
    <h3>`)
//line views/vwelcome/Welcome.html:26
	qw422016.E().S(ps.Title)
//line views/vwelcome/Welcome.html:26
	qw422016.N().S(`</h3>
    <p>It looks like you started `)
//line views/vwelcome/Welcome.html:27
	qw422016.E().S(util.AppName)
//line views/vwelcome/Welcome.html:27
	qw422016.N().S(` in a directory without a project. Set your project's details using this form and we'll get started</p>
  </div>
  <form action="" method="post">
    <div class="card">
      <h3>Create your project</h3>
      <table class="mt min-200 expanded">
        <tbody>
          `)
//line views/vwelcome/Welcome.html:34
	components.StreamTableInput(qw422016, "key", "Key", prj.Key, 5, project.Helpers["key"]...)
//line views/vwelcome/Welcome.html:34
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:35
	components.StreamTableInput(qw422016, "name", "Name", strings.TrimSuffix(prj.Name, " (missing)"), 5, project.Helpers["name"]...)
//line views/vwelcome/Welcome.html:35
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:36
	components.StreamTableInput(qw422016, "version", "Version", prj.Version, 5, project.Helpers["version"]...)
//line views/vwelcome/Welcome.html:36
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:37
	components.StreamTableInput(qw422016, "org", "Organization", info.Org, 5, project.Helpers["org"]...)
//line views/vwelcome/Welcome.html:37
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:38
	components.StreamTableInput(qw422016, "package", "Package", prj.Package, 5, project.Helpers["package"]...)
//line views/vwelcome/Welcome.html:38
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:39
	components.StreamTableInput(qw422016, "homepage", "Homepage", info.Homepage, 5, project.Helpers["homepage"]...)
//line views/vwelcome/Welcome.html:39
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:40
	components.StreamTableInput(qw422016, "sourcecode", "Source Code", info.Sourcecode, 5, project.Helpers["sourcecode"]...)
//line views/vwelcome/Welcome.html:40
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:41
	components.StreamTableInput(qw422016, "summary", "Summary", info.Summary, 5, project.Helpers["summary"]...)
//line views/vwelcome/Welcome.html:41
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:42
	components.StreamTableInputNumber(qw422016, "port", "Port", prj.Port, 5, project.Helpers["port"]...)
//line views/vwelcome/Welcome.html:42
	qw422016.N().S(`
          `)
//line views/vwelcome/Welcome.html:43
	components.StreamTableInput(qw422016, "license", "License", info.License, 5, project.Helpers["license"]...)
//line views/vwelcome/Welcome.html:43
	qw422016.N().S(`
        </tbody>
      </table>
    </div>
    <div class="card">
      <h3>Modules</h3>
      <table class="mt min-200">
        <tbody>
`)
//line views/vwelcome/Welcome.html:51
	for _, mod := range as.Services.Modules.Modules() {
//line views/vwelcome/Welcome.html:51
		qw422016.N().S(`        <tr>
          <th class="shrink">
            <label>
              <input type="checkbox" name="modules" value="`)
//line views/vwelcome/Welcome.html:55
		qw422016.E().S(mod.Key)
//line views/vwelcome/Welcome.html:55
		qw422016.N().S(`"`)
//line views/vwelcome/Welcome.html:55
		if prj.HasModule(mod.Key) {
//line views/vwelcome/Welcome.html:55
			qw422016.N().S(`checked="checked"`)
//line views/vwelcome/Welcome.html:55
		}
//line views/vwelcome/Welcome.html:55
		qw422016.N().S(`/>
              &nbsp;`)
//line views/vwelcome/Welcome.html:56
		qw422016.E().S(mod.Title())
//line views/vwelcome/Welcome.html:56
		qw422016.N().S(`
            </label>
          </th>
          <td>`)
//line views/vwelcome/Welcome.html:59
		qw422016.E().S(mod.Description)
//line views/vwelcome/Welcome.html:59
		qw422016.N().S(`</td>
        </tr>
`)
//line views/vwelcome/Welcome.html:61
	}
//line views/vwelcome/Welcome.html:61
	qw422016.N().S(`        </tbody>
      </table>
    </div>
    <div class="card">
      <button type="submit">Save</button>
      <button type="reset">Reset</button>
    </div>
  </form>
`)
//line views/vwelcome/Welcome.html:70
}

//line views/vwelcome/Welcome.html:70
func (p *Welcome) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vwelcome/Welcome.html:70
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vwelcome/Welcome.html:70
	p.StreamBody(qw422016, as, ps)
//line views/vwelcome/Welcome.html:70
	qt422016.ReleaseWriter(qw422016)
//line views/vwelcome/Welcome.html:70
}

//line views/vwelcome/Welcome.html:70
func (p *Welcome) Body(as *app.State, ps *cutil.PageState) string {
//line views/vwelcome/Welcome.html:70
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vwelcome/Welcome.html:70
	p.WriteBody(qb422016, as, ps)
//line views/vwelcome/Welcome.html:70
	qs422016 := string(qb422016.B)
//line views/vwelcome/Welcome.html:70
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vwelcome/Welcome.html:70
	return qs422016
//line views/vwelcome/Welcome.html:70
}