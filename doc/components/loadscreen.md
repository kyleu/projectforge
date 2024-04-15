# Load Screen

This component is useful as an interstitial page appearing before a long request, to warn the user it might take a while.

```go
package controllers

func LongTask(w http.ResponseWriter, r *http.Request) {
	controller.Act("long.task", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		if r.URL.Query().Get("hasloaded") != util.BoolTrue {
			cutil.URLAddQuery(r.URL, "hasloaded", util.BoolTrue)
			page := &vpage.Load{URL: r.URL.String(), Title: "Hang Tight"}
			return controller.Render(r, as, page, ps, "breadcrumb")
		}
		return "/welcome", nil
	})
}
```
