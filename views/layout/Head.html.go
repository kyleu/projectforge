// Code generated by qtc from "Head.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// <!-- Content managed by Project Forge, see [projectforge.md] for details. -->

//line views/layout/Head.html:2
package layout

//line views/layout/Head.html:2
import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/assets"
)

//line views/layout/Head.html:9
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/layout/Head.html:9
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/layout/Head.html:9
func StreamHead(qw422016 *qt422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/layout/Head.html:9
	qw422016.N().S(`
`)
//line views/layout/Head.html:10
	thm := as.Themes.Get(ps.Profile.Theme, ps.Logger)

//line views/layout/Head.html:10
	qw422016.N().S(`  <meta charset="UTF-8">
  <title>`)
//line views/layout/Head.html:12
	qw422016.E().S(ps.TitleString())
//line views/layout/Head.html:12
	qw422016.N().S(`</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover">
  `)
//line views/layout/Head.html:14
	if ps.Description != "" {
//line views/layout/Head.html:14
		qw422016.N().S(`<meta property="description" content="`)
//line views/layout/Head.html:14
		qw422016.E().S(ps.Description)
//line views/layout/Head.html:14
		qw422016.N().S(`">
  `)
//line views/layout/Head.html:15
	}
//line views/layout/Head.html:15
	qw422016.N().S(`<meta property="og:title" content="`)
//line views/layout/Head.html:15
	qw422016.E().S(ps.TitleString())
//line views/layout/Head.html:15
	qw422016.N().S(`">
  <meta property="og:type" content="website">
  <meta property="og:image" content="/assets/`)
//line views/layout/Head.html:17
	qw422016.N().U(util.AppKey)
//line views/layout/Head.html:17
	qw422016.N().S(`.svg">
  <meta property="og:locale" content="en_US">
  <meta name="theme-color" content="`)
//line views/layout/Head.html:19
	qw422016.E().S(thm.Light.NavBackground)
//line views/layout/Head.html:19
	qw422016.N().S(`" media="(prefers-color-scheme: light)">
  <meta name="theme-color" content="`)
//line views/layout/Head.html:20
	qw422016.E().S(thm.Dark.NavBackground)
//line views/layout/Head.html:20
	qw422016.N().S(`" media="(prefers-color-scheme: dark)">`)
//line views/layout/Head.html:20
	qw422016.N().S(ps.HeaderContent)
//line views/layout/Head.html:20
	qw422016.N().S(`
  <link rel="icon" href="`)
//line views/layout/Head.html:21
	qw422016.E().S(assets.URL(`logo.svg`))
//line views/layout/Head.html:21
	qw422016.N().S(`" type="image/svg+xml">
  <style>
    `)
//line views/layout/Head.html:23
	qw422016.N().S(thm.CSS(2))
//line views/layout/Head.html:23
	qw422016.N().S(`  </style>
  <link rel="stylesheet" media="screen" href="`)
//line views/layout/Head.html:24
	qw422016.E().S(assets.URL(`client.css`))
//line views/layout/Head.html:24
	qw422016.N().S(`">
  <script type="text/javascript" src="`)
//line views/layout/Head.html:25
	qw422016.E().S(assets.URL(`client.js`))
//line views/layout/Head.html:25
	qw422016.N().S(`"></script>
`)
//line views/layout/Head.html:26
}

//line views/layout/Head.html:26
func WriteHead(qq422016 qtio422016.Writer, as *app.State, ps *cutil.PageState) {
//line views/layout/Head.html:26
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/layout/Head.html:26
	StreamHead(qw422016, as, ps)
//line views/layout/Head.html:26
	qt422016.ReleaseWriter(qw422016)
//line views/layout/Head.html:26
}

//line views/layout/Head.html:26
func Head(as *app.State, ps *cutil.PageState) string {
//line views/layout/Head.html:26
	qb422016 := qt422016.AcquireByteBuffer()
//line views/layout/Head.html:26
	WriteHead(qb422016, as, ps)
//line views/layout/Head.html:26
	qs422016 := string(qb422016.B)
//line views/layout/Head.html:26
	qt422016.ReleaseByteBuffer(qb422016)
//line views/layout/Head.html:26
	return qs422016
//line views/layout/Head.html:26
}
