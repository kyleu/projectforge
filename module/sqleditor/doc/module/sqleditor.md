# [mysql]

This is a module for [Project Forge](https://projectforge.dev). It provides an UI for accessing databases.

https://github.com/kyleu/projectforge/tree/master/module/sqleditor

### License

Licensed under [CC0](https://creativecommons.org/share-your-work/public-domain/cc0)

### Usage

Add the following to your `routes.go`:

```go
	r.GET("/sql", SQLEditor)
	r.POST("/sql", SQLRun)
```

Add the following to your `menu.go`:

```go
	&menu.Item{Key: "sql", Title: "SQL Editor", Description: "Runs custom SQL", Icon: "cog", Route: "/sql"}
```

And you'll have a SQL editor!
