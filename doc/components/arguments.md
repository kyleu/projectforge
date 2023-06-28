# Arguments

This component allows you to collect a set of arguments from a web client, with validation, before sending them to their final action.

```go
package controllers

var orderArgs = cutil.Args{
	{Key: "name", Title: "Name", Description: "Pick a name"},
	{Key: "quantity", Title: "Quantity", Description: "Select a quantity", Type: "number", Default: "100"},
}

func PlaceOrder(rc *fasthttp.RequestCtx) {
	controller.Act("place.order", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		argRes := cutil.CollectArgs(rc, orderArgs)
		if argRes.HasMissing() {
			ps.Data = argRes
			msg := "Choose some options"
			return controller.Render(rc, as, &vpage.Args{URL: rc.URI().String(), Directions: msg, ArgRes: argRes}, ps, "breadcrumb")
		}
		return "/welcome", nil
	})
}
```
