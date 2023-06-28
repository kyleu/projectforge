# Load Screen

This component is useful as an interstitial page appearing before a long request, to warn the user it might take a while.

```go
package controllers

func LongTask(rc *fasthttp.RequestCtx) {
	controller.Act("long.task", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
        if string(rc.URI().QueryArgs().Peek("hasloaded")) != "true" {
            rc.URI().QueryArgs().Set("hasloaded", "true")
            page := &vpage.Load{URL: rc.URI().String(), Title: "Hang Tight"}
            return controller.Render(rc, as, page, ps, "breadcrumb")
        }
		return "/welcome", nil
	})
}
```
