package clib

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/har"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vhar"
	"{{{ .Package }}}/views/vpage"
)

func HarList(rc *fasthttp.RequestCtx) {
	controller.Act("har.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := as.Services.Har.List(ps.Logger)
		ps.SetTitleAndData("Archives", ret)
		return controller.Render(rc, as, &vhar.List{Hars: ret}, ps, "har")
	})
}

func HarDetail(rc *fasthttp.RequestCtx) {
	controller.Act("har.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", true)
		if err != nil {
			return "", err
		}
		ret, err := as.Services.Har.Load(key)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("Archive ["+key+"]", ret)
		return controller.Render(rc, as, &vhar.Detail{Har: ret}, ps, "har", ret.Key)
	})
}

func HarDelete(rc *fasthttp.RequestCtx) {
	controller.Act("har.delete", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", true)
		if err != nil {
			return "", err
		}
		err = as.Services.Har.Delete(key, ps.Logger)
		if err != nil {
			return "", err
		}
		return controller.FlashAndRedir(true, "Archive deleted", "/har", rc, ps)
	})
}

func HarTrim(rc *fasthttp.RequestCtx) {
	controller.Act("har.trim", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", true)
		if err != nil {
			return "", err
		}
		h, err := as.Services.Har.Load(key)
		if err != nil {
			return "", err
		}
		trimArgs := cutil.Args{
			{Key: "url", Title: "URL", Description: "matches against the URL (add \"*\" on either side to match wildcards)", Type: "string"},
			{Key: "mime", Title: "MIME", Description: "matches against the MIME type of the response", Type: "string", Choices: []string{"application/json"}},
		}
		argRes := cutil.CollectArgs(rc, trimArgs)
		if argRes.HasMissing() {
			url := fmt.Sprintf("%s/trim", h.WebPath())
			ps.Data = argRes
			return controller.Render(rc, as, &vpage.Args{URL: url, Directions: "Select the requests to trim", ArgRes: argRes}, ps, "har", h.Key, "Trim")
		}
		originalCount := len(h.Entries)
		h.Entries, err = h.Entries.Find(&har.Selector{URL: argRes.Values.GetStringOpt("url"), Mime: argRes.Values.GetStringOpt("mime")})
		if err != nil {
			return "", err
		}
		newCount := len(h.Entries)
		if newCount == originalCount {
			return controller.FlashAndRedir(true, "no changes needed", h.WebPath(), rc, ps)
		}
		err = as.Services.Har.Save(h)
		if err != nil {
			return "", err
		}
		msg := fmt.Sprintf("Trimmed [%d] entries from archive", originalCount-newCount)
		return controller.FlashAndRedir(true, msg, h.WebPath(), rc, ps)
	})
}

func HarUpload(rc *fasthttp.RequestCtx) {
	controller.Act("har.upload", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		mpfrm, err := rc.MultipartForm()
		if err != nil {
			return "", err
		}
		name := strings.Join(mpfrm.Value["n"], "")
		fileHeaders, ok := mpfrm.File["f"]
		if !ok {
			return "", errors.New("no file uploaded")
		}
		if len(fileHeaders) != 1 {
			return "", errors.New("invalid file uploads")
		}
		fileHeader := fileHeaders[0]
		file, err := fileHeader.Open()
		if err != nil {
			return "", err
		}
		if name == "" {
			name = fileHeader.Filename
			if !strings.HasSuffix(name, har.Ext) {
				name += har.Ext
			}
		}

		ps.Logger.Infof("Uploaded File: %+v\n", fileHeader.Filename)
		ps.Logger.Infof("File Size: %+v\n", fileHeader.Size)
		ps.Logger.Infof("MIME Header: %+v\n", fileHeader.Header)

		defer func() { _ = file.Close() }()
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return "", err
		}
		ret := &har.Wrapper{}
		err = util.FromJSON(fileBytes, ret)
		if err != nil {
			return "", errors.Wrapf(err, "error decoding file [%s]", name)
		}
		ret.Log.Key = name
		err = as.Services.Har.Save(ret.Log)
		if err != nil {
			return "", err
		}
		msg := fmt.Sprintf("Created [%s] (%s)", name, util.ByteSizeSI(fileHeader.Size))
		redir := "/har/" + name
		return controller.FlashAndRedir(true, msg, redir, rc, ps)
	})
}
