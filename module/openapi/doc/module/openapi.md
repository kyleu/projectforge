# OpenAPI

This is a module for [Project Forge](https://projectforge.dev). It embeds the Swagger UI, using your OpenAPI specification 

https://github.com/kyleu/projectforge/tree/master/module/openapi

### License 

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

By default, no controller exposes the Swagger UI. To add your own, follow this example, then add a route to it with whatever security you prefer:

```go
func OpenAPI(w http.ResponseWriter, r *http.Request) {
	controller.Act("openapi", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		u := "http://localhost:{{{ .Port }}}/assets/openapi.json"
		ps.SetTitleAndData("OpenAPI", u)
		return controller.Render(r, as, &vopenapi.OpenAPI{URL: u}, ps, "breadcrumbs")
	})
}
```
