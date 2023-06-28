# Search

If the `search` module is enabled, a search bar will be included in your pages.
You can set `ps.SearchPath = "-"` within your action to disable its appearance.

You can learn about the search engine and the types it defines by browsing `/app/lib/search`. By default, an example search function is provided after your initial project is generated.

If the `export` module is enabled, and defined models have set `search = true`, search functions will be defined and wired in `app/lib/search/generated.go`
