package cutil

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

const (
	mimeCSV   = "text/csv"
	mimeDebug = "text/plain"
	mimeHTML  = "text/html"
	mimeJSON  = "application/json"
	mimeTOML  = "application/toml"
	mimeXML   = "text/xml"
	mimeYAML  = "application/x-yaml"

	HeaderAccept                        = "Accept"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderCacheControl                  = "Cache-Control"
	HeaderContentType                   = "Content-Type"
	HeaderReferer                       = "Referer"
)

var (
	AllowedRequestHeaders  = "*"
	AllowedResponseHeaders = "*"
)

func WriteCORS(w http.ResponseWriter) {
	h := w.Header()
	setIfEmpty := func(k string, v string) {
		if x := h.Get(k); x == "" {
			h.Set(k, v)
		}
	}
	setIfEmpty(HeaderAccessControlAllowHeaders, AllowedRequestHeaders)
	setIfEmpty(HeaderAccessControlAllowMethods, "GET,POST,DELETE,PUT,PATCH,OPTIONS,HEAD")
	if x := h.Get(HeaderReferer); x == "" {
		setIfEmpty(HeaderAccessControlAllowOrigin, "*")
	} else {
		u, err := url.Parse(x)
		if err == nil {
			o := u.Hostname()
			if u.Port() != "" {
				o += ":" + u.Port()
			}
			sch := u.Scheme
			if strings.Contains(o, ".network") {
				sch = "https"
			}
			setIfEmpty(HeaderAccessControlAllowOrigin, sch+"://"+o)
		} else {
			setIfEmpty(HeaderAccessControlAllowOrigin, "*")
		}
	}
	setIfEmpty(HeaderAccessControlAllowCredentials, util.BoolTrue)
	setIfEmpty(HeaderAccessControlExposeHeaders, AllowedResponseHeaders)
}

func RespondDebug(w *WriteCounter, r *http.Request, as *app.State, filename string, ps *PageState) (string, error) {
	return RespondJSON(w, filename, RequestCtxToMap(r, as, ps))
}

func RespondCSV(w *WriteCounter, filename string, body any) (string, error) {
	b, err := util.ToCSVBytes(body)
	if err != nil {
		return "", err
	}
	return RespondMIME(filename, mimeCSV, b, w)
}

func RespondJSON(w *WriteCounter, filename string, body any) (string, error) {
	b := util.ToJSONBytes(body, true)
	return RespondMIME(filename, mimeJSON, b, w)
}

func RespondTOML(w *WriteCounter, filename string, body any) (string, error) {
	b := util.ToTOMLBytes(body)
	return RespondMIME(filename, mimeTOML, b, w)
}

type XMLResponse struct {
	Result any `xml:"result"`
}

func RespondXML(w *WriteCounter, filename string, body any) (string, error) {
	b, err := util.ToXMLBytes(body, true)
	if err != nil {
		return "", errors.Wrapf(err, "can't serialize response of type [%T] to XML", body)
	}
	return RespondMIME(filename, mimeXML, b, w)
}

func RespondYAML(w *WriteCounter, filename string, body any) (string, error) {
	b, err := yaml.Marshal(body)
	if err != nil {
		return "", err
	}
	return RespondMIME(filename, mimeYAML, b, w)
}

func RespondMIME(filename string, mime string, ba []byte, w *WriteCounter) (string, error) {
	h := w.Header()
	switch mime {
	case mimeCSV, mimeDebug, mimeHTML, mimeJSON, mimeTOML, mimeXML, mimeYAML:
		if !strings.Contains(mime, "; charset") {
			mime += "; charset=UTF-8"
		}
	}
	h.Set(HeaderContentType, mime)
	if filename != "" {
		h.Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, filename))
	}
	WriteCORS(w)
	if len(ba) == 0 {
		return "", errors.New("no bytes available to write")
	}
	if _, err := w.Write(ba); err != nil {
		return "", errors.Wrap(err, "cannot write to response")
	}

	return "", nil
}

func RespondDownload(filename string, ba []byte, w *WriteCounter) (string, error) {
	return RespondMIME(filename, "application/octet-stream", ba, w)
}

func GetContentTypes(r *http.Request) (string, string) {
	ret := r.Header.Get(HeaderAccept)
	if ret == "" {
		ret = r.Header.Get(HeaderContentType)
	}
	if idx := strings.Index(ret, ";"); idx > -1 {
		ret = ret[0:idx]
	}
	t := QueryStringString(r.URL, "t")
	switch t {
	case util.KeyDebug:
		return mimeDebug, t
	case util.KeyCSV:
		return mimeCSV, t
	case util.KeyJSON:
		return mimeJSON, t
	case util.KeyTOML:
		return mimeTOML, t
	case util.KeyXML:
		return mimeXML, t
	case util.KeyYAML, "yml":
		return mimeYAML, t
	default:
		return strings.TrimSpace(ret), t
	}
}

func GetContentType(r *http.Request) string {
	ret, _ := GetContentTypes(r)
	return ret
}

func IsContentTypeCSV(c string) bool {
	return c == mimeCSV || c == util.KeyCSV
}

func IsContentTypeDebug(c string) bool {
	return c == mimeDebug || c == util.KeyDebug
}

func IsContentTypeJSON(c string) bool {
	return c == mimeJSON || c == "text/json" || c == util.KeyJSON
}

func IsContentTypeTOML(c string) bool {
	return c == mimeTOML || c == "text/toml" || c == util.KeyTOML
}

func IsContentTypeXML(c string) bool {
	return c == "application/xml" || c == mimeXML || c == util.KeyXML
}

func IsContentTypeYAML(c string) bool {
	return c == mimeYAML || c == "text/yaml" || c == util.KeyYAML
}
