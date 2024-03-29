// Code generated by qtc from "Form.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// <!-- Content managed by Project Forge, see [projectforge.md] for details. -->

//line views/vsearch/Form.html:2
package vsearch

//line views/vsearch/Form.html:2
import (
	"strings"

	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/components"
)

//line views/vsearch/Form.html:9
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vsearch/Form.html:9
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vsearch/Form.html:9
func StreamForm(qw422016 *qt422016.Writer, act string, q string, placeholder string, currTags []string, ps *cutil.PageState) {
//line views/vsearch/Form.html:9
	qw422016.N().S(`
`)
//line views/vsearch/Form.html:10
	if placeholder == "" {
		placeholder = "Search"
	}

//line views/vsearch/Form.html:12
	qw422016.N().S(`  <form action="`)
//line views/vsearch/Form.html:13
	qw422016.E().S(act)
//line views/vsearch/Form.html:13
	qw422016.N().S(`" method="get">
`)
//line views/vsearch/Form.html:14
	if len(currTags) > 0 {
//line views/vsearch/Form.html:14
		qw422016.N().S(`    <input type="hidden" name="tags" value="`)
//line views/vsearch/Form.html:15
		qw422016.E().S(strings.Join(currTags, `,`))
//line views/vsearch/Form.html:15
		qw422016.N().S(`" />
`)
//line views/vsearch/Form.html:16
	}
//line views/vsearch/Form.html:16
	qw422016.N().S(`    <div class="right">
      <button class="right" type="submit">`)
//line views/vsearch/Form.html:18
	components.StreamSVGRef(qw422016, "search", 22, 22, `icon`, ps)
//line views/vsearch/Form.html:18
	qw422016.N().S(`</button>
      <input class="right br0" type="text" name="q" value="`)
//line views/vsearch/Form.html:19
	qw422016.E().S(q)
//line views/vsearch/Form.html:19
	qw422016.N().S(`" placeholder="`)
//line views/vsearch/Form.html:19
	qw422016.E().S(placeholder)
//line views/vsearch/Form.html:19
	qw422016.N().S(`" />
    </div>
  </form>
`)
//line views/vsearch/Form.html:22
}

//line views/vsearch/Form.html:22
func WriteForm(qq422016 qtio422016.Writer, act string, q string, placeholder string, currTags []string, ps *cutil.PageState) {
//line views/vsearch/Form.html:22
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vsearch/Form.html:22
	StreamForm(qw422016, act, q, placeholder, currTags, ps)
//line views/vsearch/Form.html:22
	qt422016.ReleaseWriter(qw422016)
//line views/vsearch/Form.html:22
}

//line views/vsearch/Form.html:22
func Form(act string, q string, placeholder string, currTags []string, ps *cutil.PageState) string {
//line views/vsearch/Form.html:22
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vsearch/Form.html:22
	WriteForm(qb422016, act, q, placeholder, currTags, ps)
//line views/vsearch/Form.html:22
	qs422016 := string(qb422016.B)
//line views/vsearch/Form.html:22
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vsearch/Form.html:22
	return qs422016
//line views/vsearch/Form.html:22
}
