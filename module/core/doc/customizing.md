# Customizing

## Startup

Your project has a full CLI interface, run `{{{ .Key }}} help` to see available options.

When the main HTTP server starts, the code in `app/controller/init.go` is run. 
It contains `initApp`, for system startup logic, and `initAppRequest`, which is run before each HTTP request. 

### Services

The main dependencies of the project are in `app/state.go`, which defines a `State` object that should almost always be in scope. 
It contains a `Services` instance which is where we'll add all our project-specific dependencies. 
You can find it in `app/services.go`.

### HTTP Controllers

- All controller actions live in `app/controller`. Normal HTTP actions should use the `controller.Act` helper method, which extracts a session and injects dependencies.
- Your method is provided an `cutil.PageState` which contains user and session information, and includes `Title` and `Data` fields, used for HTML title and data repsonses.
- Every action supports content negotiation, you can pass a `Content-Type` header, or add [`?t=json`, `?t=yaml`, or `?t=xml`], to any URL.

### CLI Actions

- Start at `/app/cmd/run.go`; to add your own CLI actions, create a new [Coral](https://github.com/muesli/coral) action in [cmd](/m/core/fs/app/cmd)
