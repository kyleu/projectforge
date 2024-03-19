package controller

import (
	"net/http"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/sandbox"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/views/vsandbox"
)

func SandboxList(w http.ResponseWriter, r *http.Request) {
	Act("sandbox.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.SetTitleAndData("Sandboxes", sandbox.AllSandboxes)
		return Render(w, r, as, &vsandbox.List{}, ps, "sandbox")
	})
}

func SandboxRun(w http.ResponseWriter, r *http.Request) {
	Act("sandbox.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(r, "key", false)
		if err != nil {
			return "", err
		}

		sb := sandbox.AllSandboxes.Get(key)
		if sb == nil {
			return ERsp("no sandbox with key [%s]", key)
		}

		ctx, span, logger := telemetry.StartSpan(ps.Context, "sandbox."+key, ps.Logger)
		defer span.Complete()

		ret, err := sb.Run(ctx, as, logger.With("sandbox", key))
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(sb.Title, ret)
		if sb.Key == "testbed" {
			return Render(w, r, as, &vsandbox.Testbed{}, ps, "sandbox", sb.Key)
		}{{{ if .HasModule "wasmclient" }}}
		if sb.Key == "wasm" {
			return Render(w, r, as, &vsandbox.WASM{}, ps, "sandbox", sb.Key)
		}{{{ end }}}
		return Render(w, r, as, &vsandbox.Run{Key: key, Title: sb.Title, Icon: sb.Icon, Result: ret}, ps, "sandbox", sb.Key)
	})
}
