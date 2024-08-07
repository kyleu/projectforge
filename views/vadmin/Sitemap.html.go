// Code generated by qtc from "Sitemap.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vadmin/Sitemap.html:1
package vadmin

//line views/vadmin/Sitemap.html:1
import (
	"slices"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cmenu"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/views/components"
	"projectforge.dev/projectforge/views/layout"
)

//line views/vadmin/Sitemap.html:12
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vadmin/Sitemap.html:12
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vadmin/Sitemap.html:12
type Sitemap struct {
	layout.Basic
}

//line views/vadmin/Sitemap.html:16
func (p *Sitemap) StreamBody(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vadmin/Sitemap.html:16
	qw422016.N().S(`
  <div class="card">
    <h3>`)
//line views/vadmin/Sitemap.html:18
	components.StreamSVGIcon(qw422016, `star`, ps)
//line views/vadmin/Sitemap.html:18
	qw422016.N().S(` Sitemap</h3>
    <div class="mt">
      `)
//line views/vadmin/Sitemap.html:20
	StreamSitemapDetail(qw422016, ps.Menu, 1, ps)
//line views/vadmin/Sitemap.html:20
	qw422016.N().S(`
    </div>
  </div>
`)
//line views/vadmin/Sitemap.html:23
}

//line views/vadmin/Sitemap.html:23
func (p *Sitemap) WriteBody(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/vadmin/Sitemap.html:23
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vadmin/Sitemap.html:23
	p.StreamBody(qw422016, as, ps)
//line views/vadmin/Sitemap.html:23
	qt422016.ReleaseWriter(qw422016)
//line views/vadmin/Sitemap.html:23
}

//line views/vadmin/Sitemap.html:23
func (p *Sitemap) Body(as *app.State, ps *cutil.PageState) string {
//line views/vadmin/Sitemap.html:23
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vadmin/Sitemap.html:23
	p.WriteBody(qb422016, as, ps)
//line views/vadmin/Sitemap.html:23
	qs422016 := string(qb422016.B)
//line views/vadmin/Sitemap.html:23
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vadmin/Sitemap.html:23
	return qs422016
//line views/vadmin/Sitemap.html:23
}

//line views/vadmin/Sitemap.html:25
func StreamSitemapDetail(qw422016 *qt422016.Writer, m menu.Items, indent int, ps *cutil.PageState) {
//line views/vadmin/Sitemap.html:26
	components.StreamIndent(qw422016, true, 1)
//line views/vadmin/Sitemap.html:26
	qw422016.N().S(`<div class="mt">`)
//line views/vadmin/Sitemap.html:28
	components.StreamIndent(qw422016, true, 2)
//line views/vadmin/Sitemap.html:28
	qw422016.N().S(`<ul class="level-0">`)
//line views/vadmin/Sitemap.html:30
	for _, i := range m {
//line views/vadmin/Sitemap.html:31
		if i.Key != "" {
//line views/vadmin/Sitemap.html:32
			streamsitemapItemDetail(qw422016, i, []string{}, ps.Breadcrumbs, 3, ps)
//line views/vadmin/Sitemap.html:33
		}
//line views/vadmin/Sitemap.html:34
	}
//line views/vadmin/Sitemap.html:35
	components.StreamIndent(qw422016, true, 2)
//line views/vadmin/Sitemap.html:35
	qw422016.N().S(`</ul>`)
//line views/vadmin/Sitemap.html:37
	components.StreamIndent(qw422016, true, 1)
//line views/vadmin/Sitemap.html:37
	qw422016.N().S(`</div>`)
//line views/vadmin/Sitemap.html:39
}

//line views/vadmin/Sitemap.html:39
func WriteSitemapDetail(qq422016 qtio422016.Writer, m menu.Items, indent int, ps *cutil.PageState) {
//line views/vadmin/Sitemap.html:39
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vadmin/Sitemap.html:39
	StreamSitemapDetail(qw422016, m, indent, ps)
//line views/vadmin/Sitemap.html:39
	qt422016.ReleaseWriter(qw422016)
//line views/vadmin/Sitemap.html:39
}

//line views/vadmin/Sitemap.html:39
func SitemapDetail(m menu.Items, indent int, ps *cutil.PageState) string {
//line views/vadmin/Sitemap.html:39
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vadmin/Sitemap.html:39
	WriteSitemapDetail(qb422016, m, indent, ps)
//line views/vadmin/Sitemap.html:39
	qs422016 := string(qb422016.B)
//line views/vadmin/Sitemap.html:39
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vadmin/Sitemap.html:39
	return qs422016
//line views/vadmin/Sitemap.html:39
}

//line views/vadmin/Sitemap.html:41
func streamsitemapItemDetail(qw422016 *qt422016.Writer, i *menu.Item, path []string, breadcrumbs cmenu.Breadcrumbs, indent int, ps *cutil.PageState) {
//line views/vadmin/Sitemap.html:42
	components.StreamIndent(qw422016, true, indent)
//line views/vadmin/Sitemap.html:42
	qw422016.N().S(`<li><div class="mts">`)
//line views/vadmin/Sitemap.html:45
	components.StreamIndent(qw422016, true, indent+1)
//line views/vadmin/Sitemap.html:45
	qw422016.N().S(`<a href="`)
//line views/vadmin/Sitemap.html:46
	qw422016.E().S(i.Route)
//line views/vadmin/Sitemap.html:46
	qw422016.N().S(`" title="`)
//line views/vadmin/Sitemap.html:46
	qw422016.E().S(i.Desc())
//line views/vadmin/Sitemap.html:46
	qw422016.N().S(`">`)
//line views/vadmin/Sitemap.html:47
	if i.Icon != "" {
//line views/vadmin/Sitemap.html:48
		components.StreamSVGRef(qw422016, i.Icon, 16, 16, "icon", ps)
//line views/vadmin/Sitemap.html:48
		qw422016.N().S(` `)
//line views/vadmin/Sitemap.html:49
	}
//line views/vadmin/Sitemap.html:50
	qw422016.E().S(i.Title)
//line views/vadmin/Sitemap.html:50
	qw422016.N().S(`</a><div><em>`)
//line views/vadmin/Sitemap.html:52
	qw422016.E().S(i.Desc())
//line views/vadmin/Sitemap.html:52
	qw422016.N().S(`</em></div>`)
//line views/vadmin/Sitemap.html:53
	if len(i.Children) > 0 {
//line views/vadmin/Sitemap.html:53
		qw422016.N().S(`<ul class="level-`)
//line views/vadmin/Sitemap.html:54
		qw422016.N().D(len(path))
//line views/vadmin/Sitemap.html:54
		qw422016.N().S(`">`)
//line views/vadmin/Sitemap.html:55
		for _, kid := range i.Children {
//line views/vadmin/Sitemap.html:56
			if kid.Key != "" {
//line views/vadmin/Sitemap.html:57
				streamsitemapItemDetail(qw422016, kid, append(slices.Clone(path), i.Key), breadcrumbs, indent+2, ps)
//line views/vadmin/Sitemap.html:58
			}
//line views/vadmin/Sitemap.html:59
		}
//line views/vadmin/Sitemap.html:59
		qw422016.N().S(`</ul>`)
//line views/vadmin/Sitemap.html:61
	}
//line views/vadmin/Sitemap.html:61
	qw422016.N().S(`</div>`)
//line views/vadmin/Sitemap.html:63
	components.StreamIndent(qw422016, true, indent)
//line views/vadmin/Sitemap.html:63
	qw422016.N().S(`</li>`)
//line views/vadmin/Sitemap.html:65
}

//line views/vadmin/Sitemap.html:65
func writesitemapItemDetail(qq422016 qtio422016.Writer, i *menu.Item, path []string, breadcrumbs cmenu.Breadcrumbs, indent int, ps *cutil.PageState) {
//line views/vadmin/Sitemap.html:65
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vadmin/Sitemap.html:65
	streamsitemapItemDetail(qw422016, i, path, breadcrumbs, indent, ps)
//line views/vadmin/Sitemap.html:65
	qt422016.ReleaseWriter(qw422016)
//line views/vadmin/Sitemap.html:65
}

//line views/vadmin/Sitemap.html:65
func sitemapItemDetail(i *menu.Item, path []string, breadcrumbs cmenu.Breadcrumbs, indent int, ps *cutil.PageState) string {
//line views/vadmin/Sitemap.html:65
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vadmin/Sitemap.html:65
	writesitemapItemDetail(qb422016, i, path, breadcrumbs, indent, ps)
//line views/vadmin/Sitemap.html:65
	qs422016 := string(qb422016.B)
//line views/vadmin/Sitemap.html:65
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vadmin/Sitemap.html:65
	return qs422016
//line views/vadmin/Sitemap.html:65
}
