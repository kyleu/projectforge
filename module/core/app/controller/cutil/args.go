package cutil

import (
	"github.com/valyala/fasthttp"
)

type Arg struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Default     string `json:"default,omitempty"`
}

type Args []*Arg

type ArgResults struct {
	Args    Args              `json:"args"`
	Values  map[string]string `json:"values"`
	Missing []string          `json:"missing,omitempty"`
}

func CollectArgs(rc *fasthttp.RequestCtx, args Args) *ArgResults {
	ret := make(map[string]string, len(args))
	var missing []string
	for _, arg := range args {
		qa := rc.URI().QueryArgs()
		if !qa.Has(arg.Key) {
			missing = append(missing, arg.Key)
		}
		ret[arg.Key] = string(qa.Peek(arg.Key))
	}
	return &ArgResults{Args: args, Values: ret, Missing: missing}
}
