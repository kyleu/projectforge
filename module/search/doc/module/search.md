# Search

The **`search`** module adds a global search bar and a `/search` results page to your Project Forge application. Results come from provider functions you register in `app/lib/search/search.go`.

## Overview

- Global search form in the main nav (core layout)
- GET `/search` results page with auto-redirect for a single match
- Provider-based search pipeline with async fan-out
- Query parsing helpers for general terms and `key:value` pairs
- Match highlighting in the results UI

## URLs and Formats

- `GET /search?q=<query>` renders HTML results by default.
- Alternate formats are available via `?format=json` (CSV/XML/YAML as supported by `controller.Render`).

## Configuration

There are no module-specific config keys. To control where the search form points or to hide it per-page:

```go
// Set a custom search endpoint (default is /search)
ps.SearchPath = "/search"

// Hide the nav search form for a handler
ps.SearchPath = "-"
```

## Usage

1. Enable the `search` module in your project configuration.
2. Register providers in `app/lib/search/search.go` under the `PF_SECTION_START(search_functions)` block.
3. Use `search.Params` helpers to parse queries and return `result.Result` entries.

Provider signature:

```go
type Provider func(
    ctx context.Context,
    params *search.Params,
    as *app.State,
    ps *cutil.PageState,
    logger util.Logger,
) (result.Results, error)
```

### Query Helpers

- `params.Parts()` splits the query on spaces.
- `params.General()` returns terms without `:`.
- `params.Keyed()` returns a `map[string]string` for `key:value` terms.

### Result Fields

`result.Result` supports:

- `Type`, `ID`, `Title`, `Icon`, `URL`
- `Matches` for highlighted text (use `result.MatchesFor`)
- `Data` to expose JSON in a modal when `ID` is set
- `HTML` for extra inline content

## Examples

A provider that supports `user:` keyed searches and highlights matches:

```go
func searchUsers(ctx context.Context, params *search.Params, as *app.State, ps *cutil.PageState, logger util.Logger) (result.Results, error) {
    keyed := params.Keyed()
    q := params.Q

    if userQuery, ok := keyed["user"]; ok {
        q = userQuery
    }

    users := loadUsersMatching(q)
    ret := make(result.Results, 0, len(users))
    for _, u := range users {
        ret = append(ret, &result.Result{
            Type:    "user",
            ID:      u.ID,
            Title:   u.Name,
            Icon:    "profile",
            URL:     "/users/" + u.ID,
            Matches: result.MatchesFor("user", u, q),
            Data:    u,
        })
    }
    return ret, nil
}
```

Register it:

```go
// $PF_SECTION_START(search_functions)$
allProviders = append(allProviders, searchUsers)
// $PF_SECTION_END(search_functions)$
```

## Dependencies

- Requires the `core` module for navigation rendering and the `/search` route.
- If the `export` module is enabled with models, additional generated search providers are appended automatically.

## CLI

This module does not add CLI commands. Use the `/search` HTTP endpoint.

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/search
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
