// $PF_HAS_MODULE(wasmclient)$
package sandbox

import (
	"context"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

var wasm = &Sandbox{Key: "wasm", Title: "WASM", Icon: "gift", Run: onWASM}

func onWASM(_ context.Context, _ *app.State, _ util.ValueMap, _ util.Logger) (any, error) {
	ret := util.ValueMap{"status": "ok"}
	return ret, nil
}
