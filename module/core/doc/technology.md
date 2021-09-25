# Technology

## Dependencies

{{{ .Name }}} relies on a ton of open source projects. First among them is the [Go language][1] itself. Other libraries include:

- [fasthttp][2]
- [quicktemplate][3]
- [chroma][4]
- [goth][5]
- [zap][6]{{{ if .HasModule "database" }}}
- [sqlx][7]{{{ end }}}{{{ if .HasModule "postgres" }}}
- [pgx][8]{{{ end }}}{{{ if .HasModule "sqlite" }}}
- [sqlite][9]{{{ end }}}


[1]: https://golang.org "What a great contribution to the world of engineering"
[2]: https://github.com/valyala/fasthttp "So much faster than the stdlib, and only slightly more annoying to work with"
[3]: https://github.com/valyala/quicktemplate "The only compile-time template engine that lets you control whitespace"
[4]: https://github.com/alecthomas/chroma "Renders a syntax-highlighted table in a surprisingly small amount of time"
[5]: https://github.com/markbates/goth "Handles OAuth for dozens of providers, works every time"
[6]: https://go.uber.org/zap "Crazy fast logging, with a custom encoder to dump tons of debug info"{{{ if .HasModule "database" }}}
[7]: https://github.com/jmoiron/sqlx "Provides enhancements to the stdlib's sql package, super handy"{{{ end }}}{{{ if .HasModule "postgres" }}}
[8]: https://github.com/jackc/pgx "Handles (most of) the crazy types that PostgreSQL supports"{{{ end }}}{{{ if .HasModule "sqlite" }}}
[9]: https://modernc.org/sqlite "A version of SQLite that was compiled to Go by a machine"{{{ end }}}
