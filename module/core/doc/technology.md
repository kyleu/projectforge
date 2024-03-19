# Technology

## Dependencies

{{{ .Name }}} relies on a ton of open source projects. First among them is the [Go language][1] itself. Other libraries include:

- [quicktemplate][3]
- [chroma][4]
- [goth][5]
- [zap][6]
- [goreleaser][7]
- [coral][8]
- [dateparse][9]
- [go-pluralize][10]
- [go-humanize][11]
- [gomarkdown][12]
- [json-iterator][13]
- [configdir][14]
- [pkg/errors][15]
- [prometheus][16]
- [lo][17]
- [opentelemetry][18]{{{ if .HasModule "database" }}}
- [sqlx][19]{{{ end }}}{{{ if .HasModule "postgres" }}}
- [pgx][20]{{{ end }}}{{{ if .HasModule "mysql" }}}
- [mysql][21]{{{ end }}}{{{ if .HasModule "sqlite" }}}
- [sqlite][22]{{{ end }}}{{{ if .HasModule "ios" }}}
- [xcodegen][23]{{{ end }}}


[1]: https://golang.org "What a great contribution to the world of engineering"
[3]: https://github.com/valyala/quicktemplate "The only compile-time template engine that lets you control whitespace"
[4]: https://github.com/alecthomas/chroma "Renders a syntax-highlighted table in a surprisingly small amount of time"
[5]: https://github.com/markbates/goth "Handles OAuth for dozens of providers, works every time"
[6]: https://go.uber.org/zap "Crazy fast logging, with a custom encoder to dump tons of debug info"
[7]: https://goreleaser.com "Builds projects in all sorts of formats"
[8]: https://github.com/muesli/coral "Provides a CLI interface without the bloat"
[9]: https://github.com/araddon/dateparse "Parses dates in all sorts of formats"
[10]: https://github.com/gertd/go-pluralize "Provides plural forms of English words"
[11]: https://github.com/dustin/go-humanize "Displays friendly relative time formats"
[12]: https://github.com/gomarkdown/markdown "Render Markdown files as HTML"
[13]: https://github.com/json-iterator/go "Fast JSON parsing and serialization"
[14]: https://github.com/kirsle/configdir "Provides access to OS-specific directories"
[15]: https://github.com/pkg/errors "Errors with stack traces and detailed logging"
[16]: https://github.com/prometheus/client_golang "Metrics for all aspects of the system"
[17]: https://github.com/samber/lo "Functional programming conveniences, used everywhere"
[18]: https://go.opentelemetry.io/otel "Telemetry for full system tracing"{{{ if .HasModule "database" }}}
[19]: https://github.com/jmoiron/sqlx "Provides enhancements to the stdlib's sql package, super handy"{{{ end }}}{{{ if .HasModule "postgres" }}}
[20]: https://github.com/jackc/pgx "Handles (most of) the crazy types that PostgreSQL supports"{{{ end }}}{{{ if .HasModule "mysql" }}}
[21]: https://github.com/go-sql-driver/mysql "The Golang MySQL driver, does what it says on the tin"{{{ end }}}{{{ if .HasModule "sqlite" }}}
[22]: https://modernc.org/sqlite "A version of SQLite that was compiled to Go by a machine"{{{ end }}}{{{ if .HasModule "ios" }}}
[23]: https://github.com/yonaskolb/XcodeGen "Generates messy iOS XCode projects"{{{ end }}}
