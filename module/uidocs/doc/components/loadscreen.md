# Load Screen

An interstitial loading screen component that provides user feedback during long-running operations. This component helps manage user expectations by displaying a loading interface before redirecting to time-consuming processes. It works completely without JavaScript, making it a great choice for expensive server-side rendering.

## Key Features

- **No JavaScript Required**: Full functionality using pure CSS and HTML
- **User Feedback**: Immediate visual confirmation that the request is being processed
- **Expectation Management**: Clear messaging about potential wait times
- **Automatic Redirect**: Seamless transition to the actual operation
- **Customizable Messaging**: Flexible title and description text

## How It Works

The load screen uses a two-step process:

1. **Initial Request**: User submits a request for a long-running operation
2. **Load Screen Display**: If the request hasn't been processed yet (no `hasloaded` parameter), show the load screen
3. **Automatic Redirect**: The load screen automatically redirects to the same URL with `hasloaded=true`
4. **Process Execution**: The second request processes the actual operation
5. **Final Redirect**: User is redirected to the completion page

This pattern prevents users from seeing blank pages or timeouts during long operations.

## Basic Usage

### Simple Load Screen

```go
package controllers

func LongTask(w http.ResponseWriter, r *http.Request) {
    controller.Act("long.task", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
        if cutil.QueryStringString(r.URL, "hasloaded") != util.BoolTrue {
            cutil.URLAddQuery(r.URL, "hasloaded", util.BoolTrue)
            page := &vpage.Load{URL: r.URL.String(), Title: "Hang Tight", Description: "Working on it..."}
            return controller.Render(r, as, page, ps, "breadcrumb")
        }

        err := performLongOperation()
        if err != nil {
            return "", err
        }
        return "/destination", nil
    })
}
```
