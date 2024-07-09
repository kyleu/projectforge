# Arguments

This component allows you to collect a set of arguments from a web client, with validation, before sending them to their final action.

```go
package controllers

var orderArgs = util.FieldDescs{
	{Key: "name", Title: "Name", Description: "Pick a name"},
	{Key: "quantity", Title: "Quantity", Description: "Select a quantity", Type: "number", Default: "100"},
}

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	controller.Act("place.order", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		argRes := util.FieldDescCollect(r, orderArgs)
		if argRes.HasMissing() {
			ps.Data = argRes
			msg := "Choose some options"
			return controller.Render(r, as, &vpage.Args{URL: r.URL.String(), Directions: msg, ArgRes: argRes}, ps, "breadcrumb")
		}
		return "/welcome", nil
	})
}
```
