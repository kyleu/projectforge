// Code generated by qtc from "Bool.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// <!-- Content managed by Project Forge, see [projectforge.md] for details. -->

//line views/components/edit/Bool.html:2
package edit

//line views/components/edit/Bool.html:2
import (
	"fmt"

	"projectforge.dev/projectforge/app/util"
)

//line views/components/edit/Bool.html:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/components/edit/Bool.html:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/components/edit/Bool.html:8
func StreamBool(qw422016 *qt422016.Writer, key string, id string, x any, nullable bool) {
//line views/components/edit/Bool.html:9
	b, isBool := x.(bool)

//line views/components/edit/Bool.html:10
	if isBool && b {
//line views/components/edit/Bool.html:10
		qw422016.N().S(`<label class="radiolabel"><input value="`)
//line views/components/edit/Bool.html:11
		qw422016.E().S(util.BoolTrue)
//line views/components/edit/Bool.html:11
		qw422016.N().S(`" name="`)
//line views/components/edit/Bool.html:11
		qw422016.E().S(key)
//line views/components/edit/Bool.html:11
		qw422016.N().S(`" type="radio" checked="checked" />`)
//line views/components/edit/Bool.html:11
		qw422016.E().S(util.BoolTrue)
//line views/components/edit/Bool.html:11
		qw422016.N().S(`</label>`)
//line views/components/edit/Bool.html:12
	} else {
//line views/components/edit/Bool.html:12
		qw422016.N().S(`<label class="radiolabel"><input value="`)
//line views/components/edit/Bool.html:13
		qw422016.E().S(util.BoolTrue)
//line views/components/edit/Bool.html:13
		qw422016.N().S(`" name="`)
//line views/components/edit/Bool.html:13
		qw422016.E().S(key)
//line views/components/edit/Bool.html:13
		qw422016.N().S(`" type="radio" />`)
//line views/components/edit/Bool.html:13
		qw422016.E().S(util.BoolTrue)
//line views/components/edit/Bool.html:13
		qw422016.N().S(`</label>`)
//line views/components/edit/Bool.html:14
	}
//line views/components/edit/Bool.html:16
	if isBool && !b {
//line views/components/edit/Bool.html:16
		qw422016.N().S(`<label class="radiolabel"><input value="`)
//line views/components/edit/Bool.html:17
		qw422016.E().S(util.BoolFalse)
//line views/components/edit/Bool.html:17
		qw422016.N().S(`" name="`)
//line views/components/edit/Bool.html:17
		qw422016.E().S(key)
//line views/components/edit/Bool.html:17
		qw422016.N().S(`" type="radio" checked="checked" />`)
//line views/components/edit/Bool.html:17
		qw422016.E().S(util.BoolFalse)
//line views/components/edit/Bool.html:17
		qw422016.N().S(`</label>`)
//line views/components/edit/Bool.html:18
	} else {
//line views/components/edit/Bool.html:18
		qw422016.N().S(`<label class="radiolabel"><input value="`)
//line views/components/edit/Bool.html:19
		qw422016.E().S(util.BoolFalse)
//line views/components/edit/Bool.html:19
		qw422016.N().S(`" name="`)
//line views/components/edit/Bool.html:19
		qw422016.E().S(key)
//line views/components/edit/Bool.html:19
		qw422016.N().S(`" type="radio" />`)
//line views/components/edit/Bool.html:19
		qw422016.E().S(util.BoolFalse)
//line views/components/edit/Bool.html:19
		qw422016.N().S(`</label>`)
//line views/components/edit/Bool.html:20
	}
//line views/components/edit/Bool.html:22
	if nullable {
//line views/components/edit/Bool.html:23
		if x == nil {
//line views/components/edit/Bool.html:23
			qw422016.N().S(`<label class="radiolabel"><input value="∅" name="`)
//line views/components/edit/Bool.html:24
			qw422016.E().S(key)
//line views/components/edit/Bool.html:24
			qw422016.N().S(`" type="radio" checked="checked" /> nil</label>`)
//line views/components/edit/Bool.html:25
		} else {
//line views/components/edit/Bool.html:25
			qw422016.N().S(`<label class="radiolabel"><input value="∅" name="`)
//line views/components/edit/Bool.html:26
			qw422016.E().S(key)
//line views/components/edit/Bool.html:26
			qw422016.N().S(`" type="radio" /> nil</label>`)
//line views/components/edit/Bool.html:27
		}
//line views/components/edit/Bool.html:28
	}
//line views/components/edit/Bool.html:29
}

//line views/components/edit/Bool.html:29
func WriteBool(qq422016 qtio422016.Writer, key string, id string, x any, nullable bool) {
//line views/components/edit/Bool.html:29
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/edit/Bool.html:29
	StreamBool(qw422016, key, id, x, nullable)
//line views/components/edit/Bool.html:29
	qt422016.ReleaseWriter(qw422016)
//line views/components/edit/Bool.html:29
}

//line views/components/edit/Bool.html:29
func Bool(key string, id string, x any, nullable bool) string {
//line views/components/edit/Bool.html:29
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/edit/Bool.html:29
	WriteBool(qb422016, key, id, x, nullable)
//line views/components/edit/Bool.html:29
	qs422016 := string(qb422016.B)
//line views/components/edit/Bool.html:29
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/edit/Bool.html:29
	return qs422016
//line views/components/edit/Bool.html:29
}

//line views/components/edit/Bool.html:31
func StreamBoolVertical(qw422016 *qt422016.Writer, key string, title string, value bool, indent int, help ...string) {
//line views/components/edit/Bool.html:32
	StreamRadioVertical(qw422016, key, title, fmt.Sprint(value), []string{util.BoolTrue, util.BoolFalse}, []string{"True", "False"}, indent)
//line views/components/edit/Bool.html:33
}

//line views/components/edit/Bool.html:33
func WriteBoolVertical(qq422016 qtio422016.Writer, key string, title string, value bool, indent int, help ...string) {
//line views/components/edit/Bool.html:33
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/edit/Bool.html:33
	StreamBoolVertical(qw422016, key, title, value, indent, help...)
//line views/components/edit/Bool.html:33
	qt422016.ReleaseWriter(qw422016)
//line views/components/edit/Bool.html:33
}

//line views/components/edit/Bool.html:33
func BoolVertical(key string, title string, value bool, indent int, help ...string) string {
//line views/components/edit/Bool.html:33
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/edit/Bool.html:33
	WriteBoolVertical(qb422016, key, title, value, indent, help...)
//line views/components/edit/Bool.html:33
	qs422016 := string(qb422016.B)
//line views/components/edit/Bool.html:33
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/edit/Bool.html:33
	return qs422016
//line views/components/edit/Bool.html:33
}

//line views/components/edit/Bool.html:35
func StreamBoolTable(qw422016 *qt422016.Writer, key string, title string, value bool, indent int, help ...string) {
//line views/components/edit/Bool.html:36
	StreamRadioTable(qw422016, key, title, fmt.Sprint(value), []string{util.BoolTrue, util.BoolFalse}, []string{"True", "False"}, indent)
//line views/components/edit/Bool.html:37
}

//line views/components/edit/Bool.html:37
func WriteBoolTable(qq422016 qtio422016.Writer, key string, title string, value bool, indent int, help ...string) {
//line views/components/edit/Bool.html:37
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/edit/Bool.html:37
	StreamBoolTable(qw422016, key, title, value, indent, help...)
//line views/components/edit/Bool.html:37
	qt422016.ReleaseWriter(qw422016)
//line views/components/edit/Bool.html:37
}

//line views/components/edit/Bool.html:37
func BoolTable(key string, title string, value bool, indent int, help ...string) string {
//line views/components/edit/Bool.html:37
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/edit/Bool.html:37
	WriteBoolTable(qb422016, key, title, value, indent, help...)
//line views/components/edit/Bool.html:37
	qs422016 := string(qb422016.B)
//line views/components/edit/Bool.html:37
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/edit/Bool.html:37
	return qs422016
//line views/components/edit/Bool.html:37
}