# Customizing your project

## Startup

Your project has a full CLI interface, run `projectforge help` to see available options. 

When the main HTTP server starts, the code in `app/controller/init.go` is run. 
It contains `initApp`, for system startup logic, and `initAppRequest`, which is run before each HTTP request. 

## Services

The main dependencies of the project are in `app/state.go`, which deines a `State` object that should almost always be in scope. 
It contains a `Services` instance which is where we'll add all our project-specific dependencies. 
You can find it in `app/services.go`.

