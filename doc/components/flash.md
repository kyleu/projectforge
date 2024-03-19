# Flash Notifications

A "flash notification" is a small message appearing on page load, and cleared before the next request. 
It's rendered as a panel in the upper-right, contains a close button (no JS required), and fades away after a few seconds (JS required).

Usually you'll want to set a flash at the end of a request, before a redirect.

```go
// FlashAndRedir(success bool, msg string, redir string, rc *fasthttp.RequestCtx, ps *cutil.PageState)
return controller.FlashAndRedir(true, msg, u, w, ps)
```

Alternately, you can add a flash manually by appending to the `Flashes` field of `cutil.PageState` in your controller action.
