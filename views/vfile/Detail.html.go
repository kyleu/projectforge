// Code generated by qtc from "Detail.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line views/vfile/Detail.html:1
package vfile

//line views/vfile/Detail.html:1
import (
	"encoding/base64"
	"strings"
	"unicode/utf8"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/components"
)

//line views/vfile/Detail.html:13
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/vfile/Detail.html:13
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/vfile/Detail.html:13
func StreamDetail(qw422016 *qt422016.Writer, path []string, b []byte, urlPrefix string, additionalLinks map[string]string, as *app.State, ps *cutil.PageState, lineNumLinkAndTitle ...string) {
//line views/vfile/Detail.html:13
	qw422016.N().S(`
`)
//line views/vfile/Detail.html:14
	if additionalLinks != nil && len(additionalLinks) > 0 {
//line views/vfile/Detail.html:14
		qw422016.N().S(`  <div class="right">
`)
//line views/vfile/Detail.html:16
		for k, v := range additionalLinks {
//line views/vfile/Detail.html:18
			var icon string
			if iconIdx := strings.Index(k, "**"); iconIdx > -1 {
				icon = k[iconIdx+2:]
				k = k[:iconIdx]
			}
			var newWindow bool
			if strings.HasPrefix(k, "!") {
				newWindow = true
				k = k[1:]
			}

//line views/vfile/Detail.html:28
			qw422016.N().S(`    <a href="`)
//line views/vfile/Detail.html:29
			qw422016.E().S(v)
//line views/vfile/Detail.html:29
			qw422016.N().S(`"`)
//line views/vfile/Detail.html:29
			if newWindow {
//line views/vfile/Detail.html:29
				qw422016.N().S(` target="_blank"`)
//line views/vfile/Detail.html:29
			}
//line views/vfile/Detail.html:29
			qw422016.N().S(`><button>`)
//line views/vfile/Detail.html:29
			if icon != "" {
//line views/vfile/Detail.html:29
				components.StreamSVGButton(qw422016, icon, ps)
//line views/vfile/Detail.html:29
				qw422016.N().S(` `)
//line views/vfile/Detail.html:29
			}
//line views/vfile/Detail.html:29
			qw422016.E().S(k)
//line views/vfile/Detail.html:29
			qw422016.N().S(`</button></a>
`)
//line views/vfile/Detail.html:30
		}
//line views/vfile/Detail.html:30
		qw422016.N().S(`  </div>
`)
//line views/vfile/Detail.html:32
	}
//line views/vfile/Detail.html:32
	qw422016.N().S(`  <h3>
`)
//line views/vfile/Detail.html:34
	for idx, p := range path {
//line views/vfile/Detail.html:34
		qw422016.N().S(`/<a href="`)
//line views/vfile/Detail.html:34
		qw422016.E().S(urlPrefix)
//line views/vfile/Detail.html:34
		qw422016.N().S(`/`)
//line views/vfile/Detail.html:34
		qw422016.E().S(util.StringFilePath(path[:idx+1]...))
//line views/vfile/Detail.html:34
		qw422016.N().S(`">`)
//line views/vfile/Detail.html:34
		qw422016.E().S(p)
//line views/vfile/Detail.html:34
		qw422016.N().S(`</a>`)
//line views/vfile/Detail.html:34
	}
//line views/vfile/Detail.html:34
	qw422016.N().S(`    <em>(`)
//line views/vfile/Detail.html:35
	qw422016.E().S(util.ByteSizeSI(int64(len(b))))
//line views/vfile/Detail.html:35
	qw422016.N().S(`)</em>
  </h3>
  <div class="mt">
`)
//line views/vfile/Detail.html:38
	if len(b) > (1024 * 128) {
//line views/vfile/Detail.html:38
		qw422016.N().S(`    <em>File is `)
//line views/vfile/Detail.html:39
		qw422016.N().D(len(b))
//line views/vfile/Detail.html:39
		qw422016.N().S(` bytes, which is too large for the file viewer</em>
`)
//line views/vfile/Detail.html:40
	} else if utf8.Valid(b) {
//line views/vfile/Detail.html:41
		out, _ := cutil.FormatFilename(string(b), path[len(path)-1], lineNumLinkAndTitle...)

//line views/vfile/Detail.html:41
		qw422016.N().S(`    `)
//line views/vfile/Detail.html:42
		qw422016.N().S(out)
//line views/vfile/Detail.html:42
		qw422016.N().S(`
`)
//line views/vfile/Detail.html:43
	} else {
//line views/vfile/Detail.html:43
		qw422016.N().S(`
`)
//line views/vfile/Detail.html:45
		if imgType := filesystem.ImageType(path...); imgType != "" {
//line views/vfile/Detail.html:45
			qw422016.N().S(`    <img alt="Image of type [`)
//line views/vfile/Detail.html:46
			qw422016.E().S(imgType)
//line views/vfile/Detail.html:46
			qw422016.N().S(`]" src="data:image/`)
//line views/vfile/Detail.html:46
			qw422016.E().S(imgType)
//line views/vfile/Detail.html:46
			qw422016.N().S(`;base64,`)
//line views/vfile/Detail.html:46
			qw422016.E().S(base64.StdEncoding.EncodeToString(b))
//line views/vfile/Detail.html:46
			qw422016.N().S(`" />
    <hr />
`)
//line views/vfile/Detail.html:48
		}
//line views/vfile/Detail.html:48
		qw422016.N().S(`
`)
//line views/vfile/Detail.html:50
		exif, err := filesystem.ExifExtract(b)

//line views/vfile/Detail.html:51
		if err == nil {
//line views/vfile/Detail.html:51
			qw422016.N().S(`    <div class="overflow full-width">
      <table>
        <thead>
          <tr>
            <th>EXIF Name</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
`)
//line views/vfile/Detail.html:61
			for k, v := range exif {
//line views/vfile/Detail.html:61
				qw422016.N().S(`          <tr>
            <td>`)
//line views/vfile/Detail.html:63
				qw422016.E().S(k)
//line views/vfile/Detail.html:63
				qw422016.N().S(`</td>
            <td>`)
//line views/vfile/Detail.html:64
				qw422016.E().V(v)
//line views/vfile/Detail.html:64
				qw422016.N().S(`</td>
          </tr>
`)
//line views/vfile/Detail.html:66
			}
//line views/vfile/Detail.html:66
			qw422016.N().S(`        </tbody>
      </table>
    </div>
`)
//line views/vfile/Detail.html:70
		} else {
//line views/vfile/Detail.html:70
			qw422016.N().S(`    <em>File is binary and contains no exif header</em>
`)
//line views/vfile/Detail.html:72
		}
//line views/vfile/Detail.html:73
	}
//line views/vfile/Detail.html:73
	qw422016.N().S(`  </div>
`)
//line views/vfile/Detail.html:75
}

//line views/vfile/Detail.html:75
func WriteDetail(qq422016 qtio422016.Writer, path []string, b []byte, urlPrefix string, additionalLinks map[string]string, as *app.State, ps *cutil.PageState, lineNumLinkAndTitle ...string) {
//line views/vfile/Detail.html:75
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/vfile/Detail.html:75
	StreamDetail(qw422016, path, b, urlPrefix, additionalLinks, as, ps, lineNumLinkAndTitle...)
//line views/vfile/Detail.html:75
	qt422016.ReleaseWriter(qw422016)
//line views/vfile/Detail.html:75
}

//line views/vfile/Detail.html:75
func Detail(path []string, b []byte, urlPrefix string, additionalLinks map[string]string, as *app.State, ps *cutil.PageState, lineNumLinkAndTitle ...string) string {
//line views/vfile/Detail.html:75
	qb422016 := qt422016.AcquireByteBuffer()
//line views/vfile/Detail.html:75
	WriteDetail(qb422016, path, b, urlPrefix, additionalLinks, as, ps, lineNumLinkAndTitle...)
//line views/vfile/Detail.html:75
	qs422016 := string(qb422016.B)
//line views/vfile/Detail.html:75
	qt422016.ReleaseByteBuffer(qb422016)
//line views/vfile/Detail.html:75
	return qs422016
//line views/vfile/Detail.html:75
}
