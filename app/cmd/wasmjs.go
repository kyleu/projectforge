//go:build js

// Content managed by Project Forge, see [projectforge.md] for details.

package cmd

import (
	"context"
	"syscall/js"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
)

const keyWASM = "wasm"

var _router fasthttp.RequestHandler

func wasmCmd() *coral.Command {
	short := "Starts the server and exposes a WebAssembly application to scripts"
	f := func(*coral.Command, []string) error { return startWASM(_flags) }
	ret := &coral.Command{Use: keyWASM, Short: short, RunE: f}
	return ret
}

func startWASM(flags *Flags) error {
	initScript()
	if err := initIfNeeded(); err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	st, r, _, err := loadServer(flags, _logger)
	if err != nil {
		return err
	}
	defer func() { _ = st.Close(context.Background(), _logger) }()
	_router = r

	select {}
}

func initScript() {
	js.Global().Set("goFetch", js.FuncOf(FetchJS))
}

func FetchJS(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return emsg("must provide exactly one argument, an HTTP request as used by fetch")
	}
	ret, err := WASMProcess(args[0])
	if err != nil {
		return emsg(err.Error())
	}
	return ret
}

func emsg(s string) string {
	return "WASM_ERROR:" + s
}

func WASMProcess(req js.Value) (js.Value, error) {
	rc := &fasthttp.RequestCtx{Request: fasthttp.Request{}, Response: fasthttp.Response{}}
	err := populateRequest(req, rc)
	if err != nil {
		return js.Null(), err
	}
	err = runRequest(rc)
	if err != nil {
		return js.Null(), err
	}
	err = debugRC(rc)
	if err != nil {
		return js.Null(), err
	}
	return createResponse(rc)
}

func populateRequest(req js.Value, rc *fasthttp.RequestCtx) error {
	url := req.Get("url").String()
	rc.Request.URI().Update(url)

	method := req.Get("method").String()
	rc.Request.Header.SetMethod(method)

	return nil
}

func runRequest(rc *fasthttp.RequestCtx) (e error) {
	if _router == nil {
		return errors.New("didn't start app through the WASM action, no active router")
	}
	defer func() {
		if x := recover(); x != nil {
			if err, ok := x.(error); ok {
				e = err
			} else {
				panic(x)
			}
		}
	}()
	_router(rc)
	return nil
}

func createResponse(rc *fasthttp.RequestCtx) (js.Value, error) {
	rspClass := js.Global().Get("Response")
	hd := js.Global().Get("Headers").New()
	for _, b := range rc.Response.Header.PeekKeys() {
		hd.Call("set", js.ValueOf(string(b)), js.ValueOf(string(rc.Response.Header.Peek(string(b)))))
	}
	hd.Call("set", js.ValueOf("Content-Type"), js.ValueOf(string(rc.Response.Header.ContentType())))
	rspBytes := rc.Response.Body()
	opts := map[string]any{"status": rc.Response.StatusCode(), "headers": hd}
	x := rspClass.New(js.ValueOf(string(rspBytes)), js.ValueOf(opts))
	return x, nil
}

func debugRC(rc *fasthttp.RequestCtx) error {
	println(util.ToJSON(cutil.RequestCtxToMap(rc, nil)))
	return nil
}
