# Technology

## Dependencies

{{{ .Name }}} relies on a ton of open source projects. First among them is the [Go language][1] itself. Other libraries include:

- [fasthttp][2]
- [quicktemplate][3]
- [chroma][4]
- [goth][5]
- [zap][6]
- [goreleaser][7]{{{ if .HasModule "database" }}}
- [sqlx][8]{{{ end }}}{{{ if .HasModule "postgres" }}}
- [pgx][9]{{{ end }}}{{{ if .HasModule "mysql" }}}
- [mysql][10]{{{ end }}}{{{ if .HasModule "sqlite" }}}
- [sqlite][11]{{{ end }}}{{{ if .HasModule "ios" }}}
- [xcodegen][12]{{{ end }}}


[1]: https://golang.org "What a great contribution to the world of engineering"
[2]: https://github.com/valyala/fasthttp "So much faster than the stdlib, and only slightly more annoying to work with"
[3]: https://github.com/valyala/quicktemplate "The only compile-time template engine that lets you control whitespace"
[4]: https://github.com/alecthomas/chroma "Renders a syntax-highlighted table in a surprisingly small amount of time"
[5]: https://github.com/markbates/goth "Handles OAuth for dozens of providers, works every time"
[6]: https://go.uber.org/zap "Crazy fast logging, with a custom encoder to dump tons of debug info"
[7]: https://goreleaser.com "Builds projects in all sorts of formats"{{{ if .HasModule "database" }}}
[8]: https://github.com/jmoiron/sqlx "Provides enhancements to the stdlib's sql package, super handy"{{{ end }}}{{{ if .HasModule "postgres" }}}
[9]: https://github.com/jackc/pgx "Handles (most of) the crazy types that PostgreSQL supports"{{{ end }}}{{{ if .HasModule "mysql" }}}
[10]: https://github.com/go-sql-driver/mysql "The Golang MySQL driver, does what it says on the tin"{{{ end }}}{{{ if .HasModule "sqlite" }}}
[11]: https://modernc.org/sqlite "A version of SQLite that was compiled to Go by a machine"{{{ end }}}{{{ if .HasModule "ios" }}}
[12]: https://github.com/yonaskolb/XcodeGen "Generates messy iOS XCode projects"{{{ end }}}
