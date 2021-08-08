package module

import (
	"fmt"
	"net/url"
)

type Resolution struct {
	Title   string            `json:"title"`
	Project string            `json:"project"`
	Action  string            `json:"action"`
	Args    map[string]string `json:"args,omitempty"`
}

func (r *Resolution) URL() string {
	ret := fmt.Sprintf("/run/%s/%s", r.Project, r.Action)
	if len(r.Args) > 0 {
		qs := make(url.Values, len(r.Args))
		for k, v := range r.Args {
			qs[k] = []string{v}
		}
		ret += "?" + qs.Encode()
	}
	return ret
}

type Resolutions []*Resolution
