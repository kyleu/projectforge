//go:build js
// +build js

package cmd

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"syscall/js"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/controller/cutil"
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

	st, r, logger, err := loadServer(flags, util.RootLogger)
	if err != nil {
		return err
	}
	logger.Infof("Started WASM server")
	defer func() { _ = st.Close(context.Background(), util.RootLogger) }()
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
	r, err := populateRequest(req, headers, reqBody)
	if err != nil {
		return js.Null(), err
	}
	rsp, err := runRequest(r)
	if err != nil {
		return js.Null(), err
	}
	return createResponse(rsp)
}

func populateRequest(req js.Value, headers js.Value, reqBody js.Value) (r *http.Request, e error) {
	defer func() {
		if x := recover(); x != nil {
			if err, ok := x.(error); ok {
				e = err
			} else {
				panic(x)
			}
		}
	}()

	ret, err := http.NewRequestWithContext(context.Background(), req.Get("method").String(), req.Get("url").String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	ret.Host = ret.URL.Hostname()

	if !headers.IsNull() && !headers.IsUndefined() {
		for i := 0; i < headers.Length(); i++ {
			entry := headers.Index(i)
			k, v := entry.Index(0).String(), entry.Index(1).String()
			ret.Header.Set(k, v)
		}
	}

	if reqBody.IsNull() || reqBody.IsUndefined() {
		return ret, nil
	}

	body := reqBody.String()
	if body == "" {
		return ret, nil
	}
	r.Body = io.NopCloser(bytes.NewReader([]byte(body)))
	return ret, nil
}

func runRequest(r *http.Request) (ret *WASMResponse, e error) {
	if _router == nil {
		return nil, errors.New("didn't start app through the WASM action, no active router")
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

	w := &WASMResponse{Headers: http.Header{}, Body: bytes.NewBuffer(nil)}
	_router.ServeHTTP(w, r)
	return w, nil
}

type WASMResponse struct {
	StatusCode int
	Headers    http.Header
	Body       *bytes.Buffer
}

func (r *WASMResponse) Header() http.Header {
	return r.Headers
}

func (r *WASMResponse) ContentType() string {
	return r.Headers.Get(cutil.HeaderContentType)
}

func (r *WASMResponse) Write(b []byte) (int, error) {
	return r.Body.Write(b)
}

func (r *WASMResponse) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
}

func createResponse(w *WASMResponse) (js.Value, error) {
	rspClass := js.Global().Get("Response")
	hd := js.Global().Get("Headers").New()
	for k, v := range w.Header() {
		for _, x := range v {
			hd.Call("set", js.ValueOf(k), js.ValueOf(x))
		}
	}
	hd.Call("set", js.ValueOf("Content-Type"), js.ValueOf(w.ContentType))
	opts := map[string]any{"status": w.StatusCode, "headers": hd}
	x := rspClass.New(js.ValueOf(w.Body.String()), js.ValueOf(opts), js.ValueOf(map[string]any{"credentials": "same-origin"}))
	return x, nil
}
