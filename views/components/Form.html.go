// Code generated by qtc from "Form.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// <!-- Content managed by Project Forge, see [projectforge.md] for details. -->

//line views/components/Form.html:2
package components

//line views/components/Form.html:2
import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vutil"
)

//line views/components/Form.html:13
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line views/components/Form.html:13
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line views/components/Form.html:13
func StreamFormInput(qw422016 *qt422016.Writer, key string, id string, value string, placeholder ...string) {
//line views/components/Form.html:14
	if id == "" {
//line views/components/Form.html:14
		qw422016.N().S(`<input name="`)
//line views/components/Form.html:15
		qw422016.E().S(key)
//line views/components/Form.html:15
		qw422016.N().S(`" value="`)
//line views/components/Form.html:15
		qw422016.E().S(value)
//line views/components/Form.html:15
		qw422016.N().S(`"`)
//line views/components/Form.html:15
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:15
		qw422016.N().S(`/>`)
//line views/components/Form.html:16
	} else {
//line views/components/Form.html:16
		qw422016.N().S(`<input id="`)
//line views/components/Form.html:17
		qw422016.E().S(id)
//line views/components/Form.html:17
		qw422016.N().S(`" name="`)
//line views/components/Form.html:17
		qw422016.E().S(key)
//line views/components/Form.html:17
		qw422016.N().S(`" value="`)
//line views/components/Form.html:17
		qw422016.E().S(value)
//line views/components/Form.html:17
		qw422016.N().S(`"`)
//line views/components/Form.html:17
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:17
		qw422016.N().S(`/>`)
//line views/components/Form.html:18
	}
//line views/components/Form.html:19
}

//line views/components/Form.html:19
func WriteFormInput(qq422016 qtio422016.Writer, key string, id string, value string, placeholder ...string) {
//line views/components/Form.html:19
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:19
	StreamFormInput(qw422016, key, id, value, placeholder...)
//line views/components/Form.html:19
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:19
}

//line views/components/Form.html:19
func FormInput(key string, id string, value string, placeholder ...string) string {
//line views/components/Form.html:19
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:19
	WriteFormInput(qb422016, key, id, value, placeholder...)
//line views/components/Form.html:19
	qs422016 := string(qb422016.B)
//line views/components/Form.html:19
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:19
	return qs422016
//line views/components/Form.html:19
}

//line views/components/Form.html:21
func StreamFormInputPassword(qw422016 *qt422016.Writer, key string, id string, value string, placeholder ...string) {
//line views/components/Form.html:22
	if id == "" {
//line views/components/Form.html:22
		qw422016.N().S(`<input name="`)
//line views/components/Form.html:23
		qw422016.E().S(key)
//line views/components/Form.html:23
		qw422016.N().S(`" type="password" value="`)
//line views/components/Form.html:23
		qw422016.E().S(value)
//line views/components/Form.html:23
		qw422016.N().S(`"`)
//line views/components/Form.html:23
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:23
		qw422016.N().S(`/>`)
//line views/components/Form.html:24
	} else {
//line views/components/Form.html:24
		qw422016.N().S(`<input id="`)
//line views/components/Form.html:25
		qw422016.E().S(id)
//line views/components/Form.html:25
		qw422016.N().S(`" name="`)
//line views/components/Form.html:25
		qw422016.E().S(key)
//line views/components/Form.html:25
		qw422016.N().S(`" type="password" value="`)
//line views/components/Form.html:25
		qw422016.E().S(value)
//line views/components/Form.html:25
		qw422016.N().S(`"`)
//line views/components/Form.html:25
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:25
		qw422016.N().S(`/>`)
//line views/components/Form.html:26
	}
//line views/components/Form.html:27
}

//line views/components/Form.html:27
func WriteFormInputPassword(qq422016 qtio422016.Writer, key string, id string, value string, placeholder ...string) {
//line views/components/Form.html:27
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:27
	StreamFormInputPassword(qw422016, key, id, value, placeholder...)
//line views/components/Form.html:27
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:27
}

//line views/components/Form.html:27
func FormInputPassword(key string, id string, value string, placeholder ...string) string {
//line views/components/Form.html:27
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:27
	WriteFormInputPassword(qb422016, key, id, value, placeholder...)
//line views/components/Form.html:27
	qs422016 := string(qb422016.B)
//line views/components/Form.html:27
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:27
	return qs422016
//line views/components/Form.html:27
}

//line views/components/Form.html:29
func StreamFormInputNumber(qw422016 *qt422016.Writer, key string, id string, value any, placeholder ...string) {
//line views/components/Form.html:30
	if id == "" {
//line views/components/Form.html:30
		qw422016.N().S(`<input name="`)
//line views/components/Form.html:31
		qw422016.E().S(key)
//line views/components/Form.html:31
		qw422016.N().S(`" type="number" value="`)
//line views/components/Form.html:31
		qw422016.E().V(value)
//line views/components/Form.html:31
		qw422016.N().S(`"`)
//line views/components/Form.html:31
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:31
		qw422016.N().S(`/>`)
//line views/components/Form.html:32
	} else {
//line views/components/Form.html:32
		qw422016.N().S(`<input id="`)
//line views/components/Form.html:33
		qw422016.E().S(id)
//line views/components/Form.html:33
		qw422016.N().S(`" name="`)
//line views/components/Form.html:33
		qw422016.E().S(key)
//line views/components/Form.html:33
		qw422016.N().S(`" type="number" value="`)
//line views/components/Form.html:33
		qw422016.E().V(value)
//line views/components/Form.html:33
		qw422016.N().S(`"`)
//line views/components/Form.html:33
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:33
		qw422016.N().S(`/>`)
//line views/components/Form.html:34
	}
//line views/components/Form.html:35
}

//line views/components/Form.html:35
func WriteFormInputNumber(qq422016 qtio422016.Writer, key string, id string, value any, placeholder ...string) {
//line views/components/Form.html:35
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:35
	StreamFormInputNumber(qw422016, key, id, value, placeholder...)
//line views/components/Form.html:35
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:35
}

//line views/components/Form.html:35
func FormInputNumber(key string, id string, value any, placeholder ...string) string {
//line views/components/Form.html:35
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:35
	WriteFormInputNumber(qb422016, key, id, value, placeholder...)
//line views/components/Form.html:35
	qs422016 := string(qb422016.B)
//line views/components/Form.html:35
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:35
	return qs422016
//line views/components/Form.html:35
}

//line views/components/Form.html:37
func StreamFormInputFloat(qw422016 *qt422016.Writer, key string, id string, value any, placeholder ...string) {
//line views/components/Form.html:38
	if id == "" {
//line views/components/Form.html:38
		qw422016.N().S(`<input name="`)
//line views/components/Form.html:39
		qw422016.E().S(key)
//line views/components/Form.html:39
		qw422016.N().S(`" type="number" value="`)
//line views/components/Form.html:39
		qw422016.E().V(value)
//line views/components/Form.html:39
		qw422016.N().S(`" step="0.001"`)
//line views/components/Form.html:39
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:39
		qw422016.N().S(`/>`)
//line views/components/Form.html:40
	} else {
//line views/components/Form.html:40
		qw422016.N().S(`<input id="`)
//line views/components/Form.html:41
		qw422016.E().S(id)
//line views/components/Form.html:41
		qw422016.N().S(`" name="`)
//line views/components/Form.html:41
		qw422016.E().S(key)
//line views/components/Form.html:41
		qw422016.N().S(`" type="number" value="`)
//line views/components/Form.html:41
		qw422016.E().V(value)
//line views/components/Form.html:41
		qw422016.N().S(`" step="0.001"`)
//line views/components/Form.html:41
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:41
		qw422016.N().S(`/>`)
//line views/components/Form.html:42
	}
//line views/components/Form.html:43
}

//line views/components/Form.html:43
func WriteFormInputFloat(qq422016 qtio422016.Writer, key string, id string, value any, placeholder ...string) {
//line views/components/Form.html:43
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:43
	StreamFormInputFloat(qw422016, key, id, value, placeholder...)
//line views/components/Form.html:43
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:43
}

//line views/components/Form.html:43
func FormInputFloat(key string, id string, value any, placeholder ...string) string {
//line views/components/Form.html:43
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:43
	WriteFormInputFloat(qb422016, key, id, value, placeholder...)
//line views/components/Form.html:43
	qs422016 := string(qb422016.B)
//line views/components/Form.html:43
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:43
	return qs422016
//line views/components/Form.html:43
}

//line views/components/Form.html:45
func StreamFormInputTimestamp(qw422016 *qt422016.Writer, key string, id string, value *time.Time, placeholder ...string) {
//line views/components/Form.html:46
	if id == "" {
//line views/components/Form.html:46
		qw422016.N().S(`<input name="`)
//line views/components/Form.html:47
		qw422016.E().S(key)
//line views/components/Form.html:47
		qw422016.N().S(`" type="text" value="`)
//line views/components/Form.html:47
		qw422016.E().S(util.TimeToJS(value))
//line views/components/Form.html:47
		qw422016.N().S(`"`)
//line views/components/Form.html:47
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:47
		qw422016.N().S(`/>`)
//line views/components/Form.html:48
	} else {
//line views/components/Form.html:48
		qw422016.N().S(`<input id="`)
//line views/components/Form.html:49
		qw422016.E().S(id)
//line views/components/Form.html:49
		qw422016.N().S(`" name="`)
//line views/components/Form.html:49
		qw422016.E().S(key)
//line views/components/Form.html:49
		qw422016.N().S(`" type="text" value="`)
//line views/components/Form.html:49
		qw422016.E().S(util.TimeToFull(value))
//line views/components/Form.html:49
		qw422016.N().S(`"`)
//line views/components/Form.html:49
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:49
		qw422016.N().S(`/>`)
//line views/components/Form.html:50
	}
//line views/components/Form.html:51
}

//line views/components/Form.html:51
func WriteFormInputTimestamp(qq422016 qtio422016.Writer, key string, id string, value *time.Time, placeholder ...string) {
//line views/components/Form.html:51
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:51
	StreamFormInputTimestamp(qw422016, key, id, value, placeholder...)
//line views/components/Form.html:51
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:51
}

//line views/components/Form.html:51
func FormInputTimestamp(key string, id string, value *time.Time, placeholder ...string) string {
//line views/components/Form.html:51
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:51
	WriteFormInputTimestamp(qb422016, key, id, value, placeholder...)
//line views/components/Form.html:51
	qs422016 := string(qb422016.B)
//line views/components/Form.html:51
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:51
	return qs422016
//line views/components/Form.html:51
}

//line views/components/Form.html:53
func StreamFormInputUUID(qw422016 *qt422016.Writer, key string, id string, value *uuid.UUID, placeholder ...string) {
//line views/components/Form.html:55
	v := ""
	if value != nil {
		v = value.String()
	}

//line views/components/Form.html:60
	StreamFormInput(qw422016, key, id, v, placeholder...)
//line views/components/Form.html:61
}

//line views/components/Form.html:61
func WriteFormInputUUID(qq422016 qtio422016.Writer, key string, id string, value *uuid.UUID, placeholder ...string) {
//line views/components/Form.html:61
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:61
	StreamFormInputUUID(qw422016, key, id, value, placeholder...)
//line views/components/Form.html:61
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:61
}

//line views/components/Form.html:61
func FormInputUUID(key string, id string, value *uuid.UUID, placeholder ...string) string {
//line views/components/Form.html:61
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:61
	WriteFormInputUUID(qb422016, key, id, value, placeholder...)
//line views/components/Form.html:61
	qs422016 := string(qb422016.B)
//line views/components/Form.html:61
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:61
	return qs422016
//line views/components/Form.html:61
}

//line views/components/Form.html:63
func StreamFormTextarea(qw422016 *qt422016.Writer, key string, id string, rows int, value string, placeholder ...string) {
//line views/components/Form.html:64
	if id == "" {
//line views/components/Form.html:64
		qw422016.N().S(`<textarea rows="`)
//line views/components/Form.html:65
		qw422016.N().D(rows)
//line views/components/Form.html:65
		qw422016.N().S(`" name="`)
//line views/components/Form.html:65
		qw422016.E().S(key)
//line views/components/Form.html:65
		qw422016.N().S(`"`)
//line views/components/Form.html:65
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:65
		qw422016.N().S(`>`)
//line views/components/Form.html:65
		qw422016.E().S(value)
//line views/components/Form.html:65
		qw422016.N().S(`</textarea>`)
//line views/components/Form.html:66
	} else {
//line views/components/Form.html:66
		qw422016.N().S(`<textarea rows="`)
//line views/components/Form.html:67
		qw422016.N().D(rows)
//line views/components/Form.html:67
		qw422016.N().S(`" id="`)
//line views/components/Form.html:67
		qw422016.E().S(id)
//line views/components/Form.html:67
		qw422016.N().S(`" name="`)
//line views/components/Form.html:67
		qw422016.E().S(key)
//line views/components/Form.html:67
		qw422016.N().S(`"`)
//line views/components/Form.html:67
		streamphFor(qw422016, placeholder)
//line views/components/Form.html:67
		qw422016.N().S(`>`)
//line views/components/Form.html:67
		qw422016.E().S(value)
//line views/components/Form.html:67
		qw422016.N().S(`</textarea>`)
//line views/components/Form.html:68
	}
//line views/components/Form.html:69
}

//line views/components/Form.html:69
func WriteFormTextarea(qq422016 qtio422016.Writer, key string, id string, rows int, value string, placeholder ...string) {
//line views/components/Form.html:69
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:69
	StreamFormTextarea(qw422016, key, id, rows, value, placeholder...)
//line views/components/Form.html:69
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:69
}

//line views/components/Form.html:69
func FormTextarea(key string, id string, rows int, value string, placeholder ...string) string {
//line views/components/Form.html:69
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:69
	WriteFormTextarea(qb422016, key, id, rows, value, placeholder...)
//line views/components/Form.html:69
	qs422016 := string(qb422016.B)
//line views/components/Form.html:69
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:69
	return qs422016
//line views/components/Form.html:69
}

//line views/components/Form.html:71
func StreamFormSelect(qw422016 *qt422016.Writer, key string, id string, value string, opts []string, titles []string, indent int) {
//line views/components/Form.html:71
	qw422016.N().S(`<select name="`)
//line views/components/Form.html:72
	qw422016.E().S(key)
//line views/components/Form.html:72
	qw422016.N().S(`"`)
//line views/components/Form.html:72
	if id == `` {
//line views/components/Form.html:72
		qw422016.N().S(` `)
//line views/components/Form.html:72
		qw422016.N().S(`id="`)
//line views/components/Form.html:72
		qw422016.E().S(id)
//line views/components/Form.html:72
		qw422016.N().S(`"`)
//line views/components/Form.html:72
	}
//line views/components/Form.html:72
	qw422016.N().S(`>`)
//line views/components/Form.html:73
	for idx, opt := range opts {
//line views/components/Form.html:75
		title := opt
		if idx < len(titles) {
			title = titles[idx]
		}

//line views/components/Form.html:80
		vutil.StreamIndent(qw422016, true, indent+1)
//line views/components/Form.html:81
		if opt == value {
//line views/components/Form.html:81
			qw422016.N().S(`<option selected="selected" value="`)
//line views/components/Form.html:82
			qw422016.E().S(opt)
//line views/components/Form.html:82
			qw422016.N().S(`">`)
//line views/components/Form.html:82
			qw422016.E().S(title)
//line views/components/Form.html:82
			qw422016.N().S(`</option>`)
//line views/components/Form.html:83
		} else {
//line views/components/Form.html:83
			qw422016.N().S(`<option value="`)
//line views/components/Form.html:84
			qw422016.E().S(opt)
//line views/components/Form.html:84
			qw422016.N().S(`">`)
//line views/components/Form.html:84
			qw422016.E().S(title)
//line views/components/Form.html:84
			qw422016.N().S(`</option>`)
//line views/components/Form.html:85
		}
//line views/components/Form.html:86
	}
//line views/components/Form.html:87
	vutil.StreamIndent(qw422016, true, indent)
//line views/components/Form.html:87
	qw422016.N().S(`</select>`)
//line views/components/Form.html:89
}

//line views/components/Form.html:89
func WriteFormSelect(qq422016 qtio422016.Writer, key string, id string, value string, opts []string, titles []string, indent int) {
//line views/components/Form.html:89
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:89
	StreamFormSelect(qw422016, key, id, value, opts, titles, indent)
//line views/components/Form.html:89
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:89
}

//line views/components/Form.html:89
func FormSelect(key string, id string, value string, opts []string, titles []string, indent int) string {
//line views/components/Form.html:89
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:89
	WriteFormSelect(qb422016, key, id, value, opts, titles, indent)
//line views/components/Form.html:89
	qs422016 := string(qb422016.B)
//line views/components/Form.html:89
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:89
	return qs422016
//line views/components/Form.html:89
}

//line views/components/Form.html:91
func StreamFormDatalist(qw422016 *qt422016.Writer, key string, id string, value string, opts []string, titles []string, indent int, placeholder ...string) {
//line views/components/Form.html:92
	vutil.StreamIndent(qw422016, true, indent)
//line views/components/Form.html:92
	qw422016.N().S(`<input id="`)
//line views/components/Form.html:93
	qw422016.E().S(id)
//line views/components/Form.html:93
	qw422016.N().S(`" list="`)
//line views/components/Form.html:93
	qw422016.E().S(id)
//line views/components/Form.html:93
	qw422016.N().S(`-list" name="`)
//line views/components/Form.html:93
	qw422016.E().S(key)
//line views/components/Form.html:93
	qw422016.N().S(`" value="`)
//line views/components/Form.html:93
	qw422016.E().S(value)
//line views/components/Form.html:93
	qw422016.N().S(`"`)
//line views/components/Form.html:93
	streamphFor(qw422016, placeholder)
//line views/components/Form.html:93
	qw422016.N().S(`/>`)
//line views/components/Form.html:94
	vutil.StreamIndent(qw422016, true, indent)
//line views/components/Form.html:94
	qw422016.N().S(`<datalist id="`)
//line views/components/Form.html:95
	qw422016.E().S(id)
//line views/components/Form.html:95
	qw422016.N().S(`-list">`)
//line views/components/Form.html:96
	for idx, opt := range opts {
//line views/components/Form.html:98
		title := opt
		if idx < len(titles) {
			title = titles[idx]
		}

//line views/components/Form.html:103
		vutil.StreamIndent(qw422016, true, indent+1)
//line views/components/Form.html:103
		qw422016.N().S(`<option value="`)
//line views/components/Form.html:104
		qw422016.E().S(opt)
//line views/components/Form.html:104
		qw422016.N().S(`">`)
//line views/components/Form.html:104
		qw422016.E().S(title)
//line views/components/Form.html:104
		qw422016.N().S(`</option>`)
//line views/components/Form.html:105
	}
//line views/components/Form.html:106
	vutil.StreamIndent(qw422016, true, indent)
//line views/components/Form.html:106
	qw422016.N().S(`</datalist>`)
//line views/components/Form.html:108
}

//line views/components/Form.html:108
func WriteFormDatalist(qq422016 qtio422016.Writer, key string, id string, value string, opts []string, titles []string, indent int, placeholder ...string) {
//line views/components/Form.html:108
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:108
	StreamFormDatalist(qw422016, key, id, value, opts, titles, indent, placeholder...)
//line views/components/Form.html:108
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:108
}

//line views/components/Form.html:108
func FormDatalist(key string, id string, value string, opts []string, titles []string, indent int, placeholder ...string) string {
//line views/components/Form.html:108
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:108
	WriteFormDatalist(qb422016, key, id, value, opts, titles, indent, placeholder...)
//line views/components/Form.html:108
	qs422016 := string(qb422016.B)
//line views/components/Form.html:108
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:108
	return qs422016
//line views/components/Form.html:108
}

//line views/components/Form.html:110
func StreamFormRadio(qw422016 *qt422016.Writer, key string, value string, opts []string, titles []string, indent int) {
//line views/components/Form.html:111
	for idx, opt := range opts {
//line views/components/Form.html:113
		title := opt
		if idx < len(titles) {
			title = titles[idx]
		}

//line views/components/Form.html:118
		vutil.StreamIndent(qw422016, true, indent)
//line views/components/Form.html:119
		if opt == value {
//line views/components/Form.html:119
			qw422016.N().S(`<label class="radio-label"><input type="radio" name="`)
//line views/components/Form.html:120
			qw422016.E().S(key)
//line views/components/Form.html:120
			qw422016.N().S(`" value="`)
//line views/components/Form.html:120
			qw422016.E().S(opt)
//line views/components/Form.html:120
			qw422016.N().S(`" checked="checked" />`)
//line views/components/Form.html:120
			qw422016.N().S(` `)
//line views/components/Form.html:120
			qw422016.E().S(title)
//line views/components/Form.html:120
			qw422016.N().S(`</label>`)
//line views/components/Form.html:121
		} else {
//line views/components/Form.html:121
			qw422016.N().S(`<label class="radio-label"><input type="radio" name="`)
//line views/components/Form.html:122
			qw422016.E().S(key)
//line views/components/Form.html:122
			qw422016.N().S(`" value="`)
//line views/components/Form.html:122
			qw422016.E().S(opt)
//line views/components/Form.html:122
			qw422016.N().S(`" />`)
//line views/components/Form.html:122
			qw422016.N().S(` `)
//line views/components/Form.html:122
			qw422016.E().S(title)
//line views/components/Form.html:122
			qw422016.N().S(`</label>`)
//line views/components/Form.html:123
		}
//line views/components/Form.html:124
	}
//line views/components/Form.html:125
}

//line views/components/Form.html:125
func WriteFormRadio(qq422016 qtio422016.Writer, key string, value string, opts []string, titles []string, indent int) {
//line views/components/Form.html:125
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:125
	StreamFormRadio(qw422016, key, value, opts, titles, indent)
//line views/components/Form.html:125
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:125
}

//line views/components/Form.html:125
func FormRadio(key string, value string, opts []string, titles []string, indent int) string {
//line views/components/Form.html:125
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:125
	WriteFormRadio(qb422016, key, value, opts, titles, indent)
//line views/components/Form.html:125
	qs422016 := string(qb422016.B)
//line views/components/Form.html:125
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:125
	return qs422016
//line views/components/Form.html:125
}

//line views/components/Form.html:127
func StreamFormCheckbox(qw422016 *qt422016.Writer, key string, values []string, opts []string, titles []string, linebreaks bool, indent int) {
//line views/components/Form.html:128
	for idx, opt := range opts {
//line views/components/Form.html:130
		title := opt
		if idx < len(titles) {
			title = titles[idx]
		}

//line views/components/Form.html:135
		vutil.StreamIndent(qw422016, true, indent)
//line views/components/Form.html:136
		if slices.Contains(values, opt) {
//line views/components/Form.html:136
			qw422016.N().S(`<label><input type="checkbox" name="`)
//line views/components/Form.html:137
			qw422016.E().S(key)
//line views/components/Form.html:137
			qw422016.N().S(`" value="`)
//line views/components/Form.html:137
			qw422016.E().S(opt)
//line views/components/Form.html:137
			qw422016.N().S(`" checked="checked" />`)
//line views/components/Form.html:137
			qw422016.N().S(` `)
//line views/components/Form.html:137
			qw422016.E().S(title)
//line views/components/Form.html:137
			qw422016.N().S(`</label>`)
//line views/components/Form.html:138
		} else {
//line views/components/Form.html:138
			qw422016.N().S(`<label><input type="checkbox" name="`)
//line views/components/Form.html:139
			qw422016.E().S(key)
//line views/components/Form.html:139
			qw422016.N().S(`" value="`)
//line views/components/Form.html:139
			qw422016.E().S(opt)
//line views/components/Form.html:139
			qw422016.N().S(`" />`)
//line views/components/Form.html:139
			qw422016.N().S(` `)
//line views/components/Form.html:139
			qw422016.E().S(title)
//line views/components/Form.html:139
			qw422016.N().S(`</label>`)
//line views/components/Form.html:140
		}
//line views/components/Form.html:141
		if slices.Contains(values, opt) {
//line views/components/Form.html:141
			qw422016.N().S(`<br />`)
//line views/components/Form.html:143
		}
//line views/components/Form.html:144
	}
//line views/components/Form.html:145
}

//line views/components/Form.html:145
func WriteFormCheckbox(qq422016 qtio422016.Writer, key string, values []string, opts []string, titles []string, linebreaks bool, indent int) {
//line views/components/Form.html:145
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:145
	StreamFormCheckbox(qw422016, key, values, opts, titles, linebreaks, indent)
//line views/components/Form.html:145
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:145
}

//line views/components/Form.html:145
func FormCheckbox(key string, values []string, opts []string, titles []string, linebreaks bool, indent int) string {
//line views/components/Form.html:145
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:145
	WriteFormCheckbox(qb422016, key, values, opts, titles, linebreaks, indent)
//line views/components/Form.html:145
	qs422016 := string(qb422016.B)
//line views/components/Form.html:145
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:145
	return qs422016
//line views/components/Form.html:145
}

//line views/components/Form.html:147
func streamphFor(qw422016 *qt422016.Writer, phs []string) {
//line views/components/Form.html:148
	if len(phs) > 0 {
//line views/components/Form.html:148
		qw422016.N().S(` `)
//line views/components/Form.html:148
		qw422016.N().S(`placeholder="`)
//line views/components/Form.html:148
		qw422016.E().S(strings.Join(phs, "; "))
//line views/components/Form.html:148
		qw422016.N().S(`"`)
//line views/components/Form.html:148
	}
//line views/components/Form.html:149
}

//line views/components/Form.html:149
func writephFor(qq422016 qtio422016.Writer, phs []string) {
//line views/components/Form.html:149
	qw422016 := qt422016.AcquireWriter(qq422016)
//line views/components/Form.html:149
	streamphFor(qw422016, phs)
//line views/components/Form.html:149
	qt422016.ReleaseWriter(qw422016)
//line views/components/Form.html:149
}

//line views/components/Form.html:149
func phFor(phs []string) string {
//line views/components/Form.html:149
	qb422016 := qt422016.AcquireByteBuffer()
//line views/components/Form.html:149
	writephFor(qb422016, phs)
//line views/components/Form.html:149
	qs422016 := string(qb422016.B)
//line views/components/Form.html:149
	qt422016.ReleaseByteBuffer(qb422016)
//line views/components/Form.html:149
	return qs422016
//line views/components/Form.html:149
}
