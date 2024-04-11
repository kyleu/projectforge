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
	mimeJSON  = "application/json"
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

func RespondDebug(w http.ResponseWriter, r *http.Request, as *app.State, filename string, ps *PageState) (string, error) {
	return RespondJSON(w, filename, RequestCtxToMap(r, as, ps))
}

func RespondCSV(w http.ResponseWriter, filename string, body any) (string, error) {
	b, err := util.ToCSVBytes(body)
	if err != nil {
		return "", err
	}
	return RespondMIME(filename, mimeCSV, util.KeyCSV, b, w)
}

func RespondJSON(w http.ResponseWriter, filename string, body any) (string, error) {
	b := util.ToJSONBytes(body, true)
	return RespondMIME(filename, mimeJSON, util.KeyJSON, b, w)
}

type XMLResponse struct {
	Result any `xml:"result"`
}

func RespondXML(w http.ResponseWriter, filename string, body any) (string, error) {
	b, err := util.ToXMLBytes(body, true)
	if err != nil {
		return "", errors.Wrapf(err, "can't serialize response of type [%T] to XML", body)
	}
	return RespondMIME(filename, mimeXML, util.KeyXML, b, w)
}

func RespondYAML(w http.ResponseWriter, filename string, body any) (string, error) {
	b, err := yaml.Marshal(body)
	if err != nil {
		return "", err
	}
	return RespondMIME(filename, mimeYAML, util.KeyYAML, b, w)
}

func RespondMIME(filename string, mime string, ext string, ba []byte, w http.ResponseWriter) (string, error) {
	h := w.Header()
	h.Set(HeaderContentType, mime+"; charset=UTF-8")
	if filename != "" {
		if !strings.HasSuffix(filename, "."+ext) {
			filename = filename + "." + ext
		}
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

func GetContentType(r *http.Request) string {
	ret := r.Header.Get(HeaderAccept)
	if ret == "" {
		ret = r.Header.Get(HeaderContentType)
	}
	if idx := strings.Index(ret, ";"); idx > -1 {
		ret = ret[0:idx]
	}
	t := r.URL.Query().Get("t")
	switch t {
	case "debug":
		return mimeDebug
	case util.KeyCSV:
		return mimeCSV
	case util.KeyJSON:
		return mimeJSON
	case util.KeyXML:
		return mimeXML
	case util.KeyYAML, "yml":
		return mimeYAML
	default:
		return strings.TrimSpace(ret)
	}
}

func IsContentTypeCSV(c string) bool {
	return c == mimeCSV
}

func IsContentTypeDebug(c string) bool {
	return c == mimeDebug
}

func IsContentTypeJSON(c string) bool {
	return c == mimeJSON || c == "text/json"
}

func IsContentTypeXML(c string) bool {
	return c == "application/xml" || c == mimeXML
}

func IsContentTypeYAML(c string) bool {
	return c == mimeYAML || c == "text/yaml"
}
