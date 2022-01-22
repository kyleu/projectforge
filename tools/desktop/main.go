// Content managed by Project Forge, see [projectforge.md] for details.
package main

import (
	"fmt"

	"github.com/webview/webview"

	"github.com/kyleu/projectforge/app/cmd"
	"github.com/kyleu/projectforge/app/util"
)

func main() {
	if err := run(true); err != nil {
		panic(err)
	}
}

func run(debug bool) error {
	port := cmd.Lib("")
	loadWebview(debug, int(port))
	return nil
}

func loadWebview(debug bool, port int) {
	w := webview.New(debug)
	defer w.Destroy()

	w.SetTitle(util.AppName)
	w.SetSize(1280, 720, webview.HintNone)
	w.Navigate(fmt.Sprintf("http://localhost:%d", port))

	w.Run()
}
