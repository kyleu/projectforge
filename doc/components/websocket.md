# WebSockets

For the server side, a controller action must be created:

```go
func MySocket(w http.ResponseWriter, r *http.Request) {
	controller.Act("my.socket", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		channel := "channel-a"
		err = as.Services.Socket.Upgrade(rc, channel, ps.Profile, ps.Logger)
		return "", err
	})
}
```

On the client side, add this to your TypeScript:

```typescript
function open() {
  console.log("socket opened");
}

// Message: { channel: string, cmd: string, param: any }
function recv(m: Message) {
  console.log("message received", m);
}

function err(svc: string, err: string) {
  console.log("error received", svc, err);
}

document.addEventListener("DOMContentLoaded", function() {
  // first argument is `debug`, which, when set to true, adds logs
  new {your project}.Socket(true, open, recv, err, "/socket");
});
```
