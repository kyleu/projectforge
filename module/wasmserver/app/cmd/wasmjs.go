//go:build js

package cmd

import (
	"context"
	"net/http"
	"syscall/js"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

const keyWASM = "wasm"

var _router http.Handler

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

	st, r, logger, err := loadServer(flags, _logger)
	if err != nil {
		return err
	}
	logger.Infof("Started WASM server")
	defer func() { _ = st.Close(context.Background(), _logger) }()
	_router = r

	select {}
}

func initScript() {
	js.Global().Set("goFetch", js.FuncOf(FetchJS))
}

func FetchJS(_ js.Value, args []js.Value) any {
	if len(args) != 3 {
		return emsg("must provide exactly three arguments, an HTTP request as used by fetch, an array of headers, and a string representation of the body")
	}
	ret, err := WASMProcess(args[0], args[1], args[2])
	if err != nil {
		return emsg(err.Error())
	}
	return ret
}

func emsg(s string) string {
	return "WASM_ERROR:" + s
}

func WASMProcess(req js.Value, headers js.Value, reqBody js.Value) (js.Value, error) {
	rc := &fasthttp.RequestCtx{Request: fasthttp.Request{}, Response: fasthttp.Response{}}
	err := populateRequest(req, headers, reqBody, rc)
	if err != nil {
		return js.Null(), err
	}
	err = runRequest(rc)
	if err != nil {
		return js.Null(), err
	}
	return createResponse(rc)
}

func populateRequest(req js.Value, headers js.Value, reqBody js.Value, rc *fasthttp.RequestCtx) (e error) {
	defer func() {
		if x := recover(); x != nil {
			if err, ok := x.(error); ok {
				e = err
			} else {
				panic(x)
			}
		}
	}()

	url := req.Get("url").String()
	rc.Request.URI().Update(url)

	rc.Request.Header.SetHost(string(rc.Request.URI().Host()))

	method := req.Get("method").String()
	rc.Request.Header.SetMethod(method)

	if !headers.IsNull() && !headers.IsUndefined() {
		for i := 0; i < headers.Length(); i++ {
			entry := headers.Index(i)
			k, v := entry.Index(0).String(), entry.Index(1).String()
			rc.Request.Header.Set(k, v)
		}
	}

	if reqBody.IsNull() || reqBody.IsUndefined() {
		return nil
	}

	body := reqBody.String()
	if body == "" {
		return nil
	}
	rc.Request.SetBody([]byte(body))
	return nil
}

func runRequest(w http.ResponseWriter, r *http.Request) (e error) {
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

func createResponse(w http.ResponseWriter, r *http.Request) (js.Value, error) {
	rspClass := js.Global().Get("Response")
	hd := js.Global().Get("Headers").New()
	for _, b := range rc.Response.Header.PeekKeys() {
		hd.Call("set", js.ValueOf(string(b)), js.ValueOf(string(rc.Response.Header.Peek(string(b)))))
	}
	hd.Call("set", js.ValueOf("Content-Type"), js.ValueOf(string(rc.Response.Header.ContentType())))
	rspBytes := rc.Response.Body()
	opts := map[string]any{"status": rc.Response.StatusCode(), "headers": hd}
	x := rspClass.New(js.ValueOf(string(rspBytes)), js.ValueOf(opts), js.ValueOf(map[string]any{"credentials": "same-origin"}))
	return x, nil
}
