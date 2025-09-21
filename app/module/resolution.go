package module

import (
	"fmt"
	"net/url"

	"github.com/samber/lo"
)

type Resolution struct {
	Title   string            `json:"title"`
	Project string            `json:"project"`
	Action  string            `json:"action"`
	Args    map[string]string `json:"args,omitzero"`
}

func (r *Resolution) URL() string {
	ret := fmt.Sprintf("/run/%s/%s", r.Project, r.Action)
	if len(r.Args) > 0 {
		var qs url.Values = lo.MapValues(r.Args, func(v string, _ string) []string {
			return []string{v}
		})
		ret += "?" + qs.Encode()
	}
	return ret
}

type Resolutions []*Resolution
