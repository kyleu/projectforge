package cutil

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/util"
)

const (
	mimeJSON = "application/json"
	mimeXML  = "text/xml"
)

func WriteCORS(rc *fasthttp.RequestCtx) {
	rc.Response.Header.Set("Access-Control-Allow-Headers", "*")
	rc.Response.Header.Set("Access-Control-Allow-Method", "GET,POST,DELETE,PUT,PATCH,OPTIONS,HEAD")
	rc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	rc.Response.Header.Set("Access-Control-Allow-Credentials", "true")
}

func RespondDebug(rc *fasthttp.RequestCtx, filename string, body interface{}) (string, error) {
	return RespondJSON(rc, filename, requestCtxToMap(rc, body))
}

func RespondJSON(rc *fasthttp.RequestCtx, filename string, body interface{}) (string, error) {
	b := util.ToJSONBytes(body, true)
	return RespondMIME(filename, mimeJSON, "json", b, rc)
}

type XMLResponse struct {
	Result interface{} `xml:"result"`
}

func RespondXML(rc *fasthttp.RequestCtx, filename string, body interface{}) (string, error) {
	body = XMLResponse{Result: body}
	b, err := xml.Marshal(body)
	if err != nil {
		return "", errors.Wrapf(err, "can't serialize response of type [%T] to XML", body)
	}
	return RespondMIME(filename, mimeXML, "xml", b, rc)
}

func RespondMIME(filename string, mime string, ext string, ba []byte, rc *fasthttp.RequestCtx) (string, error) {
	rc.Response.Header.SetContentType(mime + "; charset=UTF-8")
	if len(filename) > 0 {
		if !strings.HasSuffix(filename, "."+ext) {
			filename = filename + "." + ext
		}
		rc.Response.Header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, filename))
	}
	WriteCORS(rc)
	if len(ba) == 0 {
		return "", errors.New("no bytes available to write")
	}
	if _, err := rc.Write(ba); err != nil {
		return "", errors.Wrap(err, "cannot write to response")
	}

	return "", nil
}

func GetContentType(rc *fasthttp.RequestCtx) string {
	ret := string(rc.Request.Header.ContentType())
	if idx := strings.Index(ret, ";"); idx > -1 {
		ret = ret[0:idx]
	}
	t := string(rc.URI().QueryArgs().Peek("t"))
	switch t {
	case "debug":
		return t
	case "json":
		return mimeJSON
	case "xml":
		return mimeXML
	default:
		return strings.TrimSpace(ret)
	}
}

func IsContentTypeJSON(c string) bool {
	return c == "application/json" || c == "text/json"
}

func IsContentTypeXML(c string) bool {
	return c == "application/xml" || c == mimeXML
}
