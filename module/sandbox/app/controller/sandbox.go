package controller

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/sandbox"
	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/views/vsandbox"
)

func SandboxList(rc *fasthttp.RequestCtx) {
	act("sandbox.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Sandboxes"
		ps.Data = sandbox.AllSandboxes
		return render(rc, as, &vsandbox.List{}, ps, "sandbox")
	})
}

func SandboxRun(rc *fasthttp.RequestCtx) {
	act("sandbox.run", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}

		sb := sandbox.AllSandboxes.Get(key)
		if sb == nil {
			return ersp("no sandbox with key [%s]", key)
		}

		ctx, span, logger := telemetry.StartSpan(ps.Context, "sandbox."+key, ps.Logger)
		defer span.Complete()

		ret, err := sb.Run(ctx, as, logger.With(zap.String("sandbox", key)))
		if err != nil {
			return "", err
		}
		ps.Title = sb.Title
		ps.Data = ret
		if sb.Key == "testbed" {
			return render(rc, as, &vsandbox.Testbed{}, ps, "sandbox", sb.Key)
		}
		return render(rc, as, &vsandbox.Run{Key: key, Title: sb.Title, Icon: sb.Icon, Result: ret}, ps, "sandbox", sb.Key)
	})
}
